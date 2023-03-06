package smc

import (
	"fmt"

	"github.com/ebfe/scard"
)

type mifareReader struct {
	card    *scard.Card
	respCmd []byte
	uid     string
}

func NewMifareReader(card *scard.Card, respCmd []byte) *mifareReader {
	return &mifareReader{
		card,
		respCmd,
		"",
	}
}

func (r *mifareReader) Select() error {
	// Send command APDU
	cmd := []byte{0xFF, 0xCA, 0x00, 0x00, 0x00}
	resp, err := r.card.Transmit(cmd)
	if err != nil {
		fmt.Printf("Failed to retrieve UID: %s\n", err)
		return err
	}
	var str string
	for _, b := range resp[1:5] {
		str += fmt.Sprintf("%d", b)
	}
	r.uid = str
	return nil
}
