package smc

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ebfe/scard"
	"github.com/somprasongd/go-thai-smartcard/pkg/logger"
	"github.com/somprasongd/go-thai-smartcard/pkg/model"
	"github.com/somprasongd/go-thai-smartcard/pkg/util"
)

type Options struct {
	ShowFaceImage bool
	ShowNhsoData  bool
	ShowLaserData bool
}

type smartCard struct {
}

func NewSmartCard() *smartCard {
	logger.LOGGER().Info("NewSmartCard")
	return &smartCard{}
}

func (s *smartCard) ListReaders() ([]string, error) {
	logger.LOGGER().Info("ListReaders")
	// Establish a context
	ctx, err := util.EstablishContext()
	if err != nil {
		return nil, err
	}
	defer util.ReleaseContext(ctx)

	// List available readers
	return util.ListReaders(ctx)
}

// func (s *smartCard) Read(readerName *string, opts *Options) (*model.Data, error) {
// 	logger.LOGGER().Info("Read")
// 	if opts == nil {
// 		opts = &Options{
// 			ShowFaceImage: true,
// 			ShowNhsoData:  false,
// 			ShowLaserData: false,
// 		}
// 	}

// 	readers := []string{}

// 	if readerName == nil {
// 		r, err := s.ListReaders()
// 		if err != nil {
// 			return nil, err
// 		}
// 		readers = r
// 	} else {
// 		readers = append(readers, *readerName)
// 	}

// 	if len(readers) == 0 {
// 		return nil, errors.New("not available readers")
// 	}

// 	// Establish a context
// 	ctx, err := util.EstablishContext()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer util.ReleaseContext(ctx)

// 	rs := util.InitReaderStates(readers)

// 	log.Println("Waiting for a Card Inserted")
// 	index, err := util.WaitUntilCardPresent(ctx, rs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	reader := readers[index]
// 	card, data, err := s.readCard(ctx, reader, opts)
// 	defer util.DisconnectCard(card)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return data, nil
// }

func (s *smartCard) readCard(ctx *scard.Context, reader string, opts *Options) (*scard.Card, *model.Data, error) {
	logger.LOGGER().Info("readCard")
	log.Printf("Connecting to card with %s", reader)
	card, err := util.ConnectCard(ctx, reader)
	if err != nil {
		log.Printf("connecting card error %s", err.Error())
		return card, nil, err
	}

	// defer func(card *scard.Card) {
	// 	if rcv := recover(); rcv != nil {
	// 		_, e := card.Status()
	// 		if e != nil {
	// 			log.Println("Recover readCard:", e.Error())
	// 			return
	// 		}
	// 		log.Println("Recover readCard:", rcv)
	// 	}
	// }(card)

	status, err := card.Status()
	if err != nil {
		log.Printf("get card status error %s", err.Error())
		return card, nil, err
	}

	cmd := util.GetResponseCommand(status.Atr)
	data := model.Data{}
	if strings.Contains(status.Reader, MIFARE_1) || strings.Contains(status.Reader, MIFARE_2) {
		mifareReader := NewMifareReader(card, cmd)
		mifareReader.Select()
		data.Personal = &model.Personal{
			Cid: mifareReader.uid,
		}
	} else if strings.Contains(status.Reader, THAI_SMC_1) || strings.Contains(status.Reader, THAI_SMC_2) {
		personalReader := NewPersonalReader(card, cmd)
		personalReader.Select()
		data.Personal = personalReader.Read(opts.ShowFaceImage)

		if opts.ShowLaserData {
			cardReader := NewCardReader(card, cmd)
			cardReader.Select()
			data.Card = &model.Card{LaserId: cardReader.ReadLaserId()}
		}

		if opts.ShowNhsoData {
			nhsoReader := NewNhsoReader(card, cmd)
			nhsoReader.Select()
			data.Nhso = nhsoReader.Read()
		}
	}
	return card, &data, nil
}

func (s *smartCard) StartDaemon(broadcast chan model.Message, opts *Options) error {
	logger.LOGGER().Info("StartDaemon")
	if opts == nil {
		opts = &Options{
			ShowFaceImage: true,
			ShowNhsoData:  false,
			ShowLaserData: false,
		}
	}

	util.InitCardReaderList()

	// Establish a context
	ctx, err := util.EstablishContext()
	if err != nil {
		log.Printf("establish context error %s\n", err.Error())
		return err
	}
	defer util.ReleaseContext(ctx)

	// chWaitReaders := make(chan []string)
	insertedCardChan := make(chan string)
	connectedCardReaders := make(chan []scard.ReaderState)
	var readers []string
	var Mtx sync.Mutex
	go func() {
		// logger.LOGGER().Info("go func chWaitReaders")
		cardReaderAmount := 0
		for {
			// fmt.Println("CURRENT ActiveCardReader : ", util.GetActiveCardReader())
			// logger.LOGGER().Info("latest connected card reader : ", cardReaderAmount)
			// List available readers
			Mtx.Lock()
			readers, err = util.ListReaders(ctx)
			if cardReaderAmount == 0 {
				util.InitCardReaderList()
			}
			if len(readers) != cardReaderAmount {
				logger.LOGGER().Info("UPDATE WATCHER")
				rs := util.InitReaderStates(readers)
				connectedCardReaders <- rs
			}
			cardReaderAmount = len(readers)
			Mtx.Unlock()
			if err != nil {
				if broadcast != nil {
					message := model.Message{
						Event: "smc-error",
						Payload: map[string]string{
							"message": err.Error(),
						},
					}
					broadcast <- message
				}
				logger.LOGGER().Error("Cannot find a smart card reader, Wait 2 seconds")
				time.Sleep(2 * time.Second)
				continue
			}
			// logger.LOGGER().Info(fmt.Sprintf("Available %d readers:\n", len(readers)))
			// for i, reader := range readers {
			// 	logger.LOGGER().Info(fmt.Sprintf("[%d] %s\n", i, reader))
			// }

			if len(readers) == 0 {
				if broadcast != nil {
					message := model.Message{
						Event: "smc-error",
						Payload: map[string]string{
							"message": "not available readers",
						},
					}
					broadcast <- message
				}
				logger.LOGGER().Error("Cannot find a smart card reader, Wait 2 seconds")
				time.Sleep(2 * time.Second)
				continue
			}
			// time.Sleep(3 * time.Second)
		}
	}()
	// readers := <-chWaitReaders
	// readers, err = util.ListReaders(ctx)
	// if err != nil {
	// 	logger.LOGGER().Error("ERROR ListReaders : ", err)
	// }

	go func() {
		for {
			newCardReaders := <-connectedCardReaders
			// logger.LOGGER().Warn("NEW CONNECTING CARD READER : ", newCardReader)
			util.AddCardReader(newCardReaders)
			go util.CardReaderWatcher(ctx, newCardReaders, insertedCardChan)
		}
	}()

	go func() {
		for {
			newInserted := <-insertedCardChan
			var card *scard.Card
			if newInserted != "" {
				logger.LOGGER().Warn("NEW INSERT : ", newInserted)
				newCard, data, err := s.readCard(ctx, newInserted, opts)
				card = newCard
				if err != nil {
					logger.LOGGER().Warn("ERROR FROM READCARD : ", err)
					util.DisconnectCard(card)
				}
				if data != nil {
					logger.LOGGER().Warn("NEW DATA : ", data)
					logger.LOGGER().Warn("FROM : ", newInserted)
					message := model.Message{
						Reader:  newInserted,
						Event:   "smc-data",
						Payload: data,
					}
					// if newInserted == MIFARE_1 || newInserted == MIFARE_2 {
					if strings.Contains(newInserted, MIFARE_1) || strings.Contains(newInserted, MIFARE_2) {
						message.Event = "mifare-data"
					}
					broadcast <- message
					util.DisconnectCard(card)
				}
			} else {
				util.DisconnectCard(card)
				logger.LOGGER().Warn("CARD WAS REMOVE")
			}
		}
	}()
	for {

	}
}
