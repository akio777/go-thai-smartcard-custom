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
	fmt.Printf("Card UID: ")
	for _, b := range resp[0:4] {
		fmt.Printf("%02X ", b)
	}
	str := fmt.Sprintf("%010d", binary.LittleEndian.Uint32(resp[0:4]))
	fmt.Printf("\nID: %s\n", str)

	r.uid = str
	return nil
}
