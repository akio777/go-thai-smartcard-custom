package smc

import (
	"encoding/binary"
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
	str = fmt.Sprintf("%d", binary.LittleEndian.Uint32(resp[0:4]))
	if len(str) < 10 {
		str = fmt.Sprintf("%d%s", 0, str)
	}
	r.uid = str
	return nil
}
