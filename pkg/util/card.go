package util

import (
	"bytes"
	"errors"
	"log"
	"strings"

	"github.com/ebfe/scard"
	"github.com/varokas/tis620"
)

func EstablishContext() (*scard.Context, error) {
	return scard.EstablishContext()
}

func ReleaseContext(ctx *scard.Context) {
	ctx.Release()
}

func ListReaders(ctx *scard.Context) ([]string, error) {
	return ctx.ListReaders()
}

func InitReaderStates(readers []string) []scard.ReaderState {
	rs := make([]scard.ReaderState, len(readers))
	for i := range rs {
		rs[i].Reader = readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}
	return rs
}

func WaitUntilCardPresent(ctx *scard.Context, rs []scard.ReaderState, insertedCardChan chan<- string) (int, error) {
	for {
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
		for i := range rs {

			rs[i].CurrentState = rs[i].EventState
			if rs[i].EventState&scard.StatePresent != 0 {
				log.Println("Card inserted")
				insertedCardChan <- rs[i].Reader
				return i, nil
			}
		}
	}
}

func WaitUntilCardRemove(ctx *scard.Context, rs []scard.ReaderState, insertedCardChan chan<- string) (int, error) {
	for {
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
		for i := range rs {

			rs[i].CurrentState = rs[i].EventState
			if rs[i].EventState&scard.StateEmpty != 0 {
				log.Println("Card removed")
				insertedCardChan <- ""
				return i, nil
			}

		}
	}
}

func ConnectCard(ctx *scard.Context, reader string) (*scard.Card, error) {
	return ctx.Connect(reader, scard.ShareExclusive, scard.ProtocolAny)
}

func DisconnectCard(card *scard.Card) error {
	if card == nil {
		return errors.New("card is nil")
	}
	return card.Disconnect(scard.UnpowerCard)
}

func GetResponseCommand(atr []byte) []byte {
	if atr[0] == 0x3B && atr[1] == 0x67 {
		return []byte{0x00, 0xc0, 0x00, 0x01}
	}
	return []byte{0x00, 0xc0, 0x00, 0x00}
}

func ReadData(card *scard.Card, cmd []byte, cmdGetResponse []byte) (string, error) {
	return readDataToString(card, cmd, cmdGetResponse, false)
}

func ReadDataThai(card *scard.Card, cmd []byte, cmdGetResponse []byte) (string, error) {
	return readDataToString(card, cmd, cmdGetResponse, true)
}

func readDataToString(card *scard.Card, cmd []byte, cmdGetResponse []byte, isTIS620 bool) (string, error) {
	_, err := card.Status()
	if err != nil {
		return "", err
	}
	// Send command APDU
	_, err = card.Transmit(cmd)
	if err != nil {
		// log.Println("Error Transmit:", err)
		return "", err
	}
	// log.Println(rsp)

	// get respond command
	cmd_respond := append(cmdGetResponse[:], cmd[len(cmd)-1])
	rsp, err := card.Transmit(cmd_respond)
	if err != nil {
		// log.Println("Error Transmit:", err)
		return "", err
	}
	// log.Println(rsp)

	if isTIS620 {
		rsp = tis620.ToUTF8(rsp)
	}

	// for i := 0; i < len(rsp)-2; i++ {
	// 	cid += fmt.Sprintf("%c", rsp[i])
	// }
	return strings.TrimSpace(string(rsp[:len(rsp)-2])), nil
}

func ReadLaserData(card *scard.Card, cmd []byte, cmdGetResponse []byte) (string, error) {
	_, err := card.Status()
	if err != nil {
		return "", err
	}
	// Send command APDU
	_, err = card.Transmit(cmd)
	if err != nil {
		return "", err
	}

	// get respond command
	cmd_respond := append(cmdGetResponse[:], 0x10) // Java use 0x80
	rsp, err := card.Transmit(cmd_respond)
	if err != nil {
		return "", err
	}
	// fmt.Printf("%v, %s, %v\n", rsp, string(bytes.Trim(rsp[:len(rsp)-2], "\x00")), len(string(bytes.Trim(rsp[:len(rsp)-2], "\x00"))))
	return strings.TrimSpace(string(bytes.Trim(rsp[:len(rsp)-2], "\x00"))), nil
}
