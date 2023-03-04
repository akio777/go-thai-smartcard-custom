package util

import (
	"fmt"
	"sync"

	"github.com/ebfe/scard"
)

var (
	Mtx              sync.Mutex
	ActiveCardReader map[string]string
)

func InitCardReaderList() {
	ActiveCardReader = make(map[string]string)
}

func GetActiveCardReader() map[string]string {
	return ActiveCardReader
}

const (
	INIT   = "init"
	EXISTS = "exists"
	DESTR  = "destr"
)

func CardReaderWatcher(ctx *scard.Context, readers []scard.ReaderState, insertedCardChan chan<- string) {
	for _, reader := range readers {
		Mtx.Lock()
		fmt.Println("reader : ", reader.Reader, ActiveCardReader[reader.Reader])
		status := ActiveCardReader[reader.Reader]
		Mtx.Unlock()
		if status == INIT {
			go func(reader scard.ReaderState) {
				Mtx.Lock()
				ActiveCardReader[reader.Reader] = EXISTS
				Mtx.Unlock()
				for {
					Mtx.Lock()
					status = ActiveCardReader[reader.Reader]
					Mtx.Unlock()
					if status == EXISTS {
						fmt.Println("WATCHING : ", reader.Reader)
						WaitUntilCardPresent(ctx, []scard.ReaderState{reader}, insertedCardChan)
						WaitUntilCardRemove(ctx, []scard.ReaderState{reader}, insertedCardChan)
					} else if status == DESTR {
						fmt.Println("REMOVE : ", reader.Reader)
						break
					}
				}
			}(reader)
		}
		// }
	}
}

func AddCardReader(newCardReaders []scard.ReaderState) {
	Mtx.Lock()
	currentReaders := len(ActiveCardReader)
	if currentReaders == 0 {
		for _, reader := range newCardReaders {
			ActiveCardReader[reader.Reader] = INIT
		}
	} else {
		fmt.Println("OLD FALSE , NEW TRUE")
		fmt.Println(ActiveCardReader)
		fmt.Println(currentReaders)
		oldCardReaders := ActiveCardReader
		for name := range oldCardReaders {
			ActiveCardReader[name] = DESTR
		}
		for _, reader := range newCardReaders {
			ActiveCardReader[reader.Reader] = EXISTS
		}
	}
	Mtx.Unlock()
}
