package orderbook

import (
	"fmt"
	"os"
)

func sideAwareIndex(i int, side byte, ladder Ladder) int {
	if side == 'S' {
		return i
	}
	return (len(*ladder.Depth)-1) - i
}

func ProcessMessage(n int, msg Message, book map[string]Ladder) bool {
	ticker := string(msg.Order.Symbol[:])
	side := string(msg.Order.Side)

	// Initialise
	for _, s := range []string{"B", "S"} {
		_, ok := book[ticker+s]
		if !ok {
			var l Ladder
			d := make([]PriceVol, 0, 10)
			l.Depth = &d
			l.Orders = map[uint64]PriceVol{}
			book[ticker+s] = l
		} else {
			break // If one is already created then both must've been
		}
	}
	ladder := book[ticker+side]

	var idx int
	switch string(msg.Order.MsgType) {
	case "A":
		if msg.Size == 0 {
			return false
		}
		idx = ladder.AddOrder(msg.Order.OrderId, msg.Price, msg.Size)
	case "U":
		idx = ladder.UpdateOrder(msg.Order.OrderId, msg.Price, msg.Size, msg.Order.Side)
	case "D":
		idx = ladder.DeleteOrder(msg.Order.OrderId)
		if idx < 0 {
			return false // Handling for deleting orders which had 0 size
		}
	case "E":
		idx = ladder.ExecuteOrder(msg.Order.OrderId, msg.Size)
	default:
		fmt.Println("Unable to process order type")
		os.Exit(1)
		return false // necessary to stop compiler from whining
	}
	return sideAwareIndex(idx, msg.Order.Side, ladder) < n
}

func FormatLadder(n int, ticker string, seqNo uint32, buySide []PriceVol, sellSide []PriceVol) string {
	str := fmt.Sprintf("%d, %s, [", seqNo, ticker)

	buyLength := len(buySide)
	for i := 1; i <= n && i <= buyLength; i++ {
		if i > 1 {
			str += ", "
		}
		pv := buySide[buyLength-i]
		str += fmt.Sprintf("(%d, %d)", pv.Price, pv.Volume)
	}

	str += "], ["

	sellLength := len(sellSide)
	for i := 0; i < n && i < sellLength; i++ {
		if i > 0 {
			str += ", "
		}
		pv := sellSide[i]
		str += fmt.Sprintf("(%d, %d)", pv.Price, pv.Volume)
	}
	str += "]"

	return str
}
