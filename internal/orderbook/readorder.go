package orderbook

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

type Header struct {
	SeqNo   uint32
	MsgSize uint32
}

type Order struct {
	MsgType byte
	Symbol  [3]byte
	OrderId uint64
	Side    byte
}

type Message struct {
	Header Header
	Order  Order
	Size   uint64
	Price  int32
}

func exit(err error, msg string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(msg)
		os.Exit(1)
	}
}

func ReadMessage(reader *bufio.Reader) (Message, int) {
	var err error
	var msg Message

	err = binary.Read(reader, binary.LittleEndian, &msg.Header)
	if err != nil && errors.Is(err, io.EOF) {
		return msg, 1
	}
	exit(err, "Failed to read package header")

	err = binary.Read(reader, binary.LittleEndian, &msg.Order)
	exit(err, "Failed to read order type")

	for i := 0; i < 3; i++ {
		reader.ReadByte()
	}

	if msg.Header.MsgSize == 16 {
		return msg, 0
	}

	err = binary.Read(reader, binary.LittleEndian, &msg.Size)
	exit(err, "Failed to read order size")
	if msg.Header.MsgSize == 24 {
		return msg, 0
	}

	err = binary.Read(reader, binary.LittleEndian, &msg.Price)
	for i := 0; i < 4; i++ {
		reader.ReadByte()
	}
	return msg, 0
}

func ReadStream(stdin *os.File) {
	reader := bufio.NewReader(stdin)

	var msg Message
	EOF := 0

	for {
		msg, EOF = ReadMessage(reader)
		if EOF == 1 {
			break
		}
		fmt.Println(msg)
	}
}
