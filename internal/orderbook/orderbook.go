package orderbook

import (
	"fmt"
	"os"
)

func ProcessMessage(n int, msg Message, book *map[string]Ladder) bool {
	ticker := string(msg.Order.Symbol[:])
	side := string(msg.Order.Side)

	// Initialise
	for _, s := range []string{"B", "S"} {
		_, ok := (*book)[ticker+s]
		if !ok {
			var l Ladder
			l.Depth = &[]priceVol{}
			l.Orders = map[uint64]priceVol{}
			(*book)[ticker+s] = l
		} else {
			break // If one is created then both are
		}
	}
	ladder := (*book)[ticker+side]

	switch string(msg.Order.MsgType) {
	case "A":
		return addOrder(msg.Order.OrderId, msg.Price, msg.Size, &ladder) < n
	case "U":
		return updateOrder(msg.Order.OrderId, msg.Price, msg.Size, &ladder, msg.Order.Side) < n
	case "D":
		idx := deleteOrder(msg.Order.OrderId, &ladder)  // Handling for deleting orders which had 0 size
		return 0 < idx && idx < n
	case "E":
		return executeOrder(msg.Order.OrderId, msg.Size, &ladder) < n
	default:
		fmt.Println("Unable to process order type")
		os.Exit(1)
		return false // necessary to stop compiler from whining
	}
}

func FormatLadder(n int, ticker string, seqNo uint32, buySide []priceVol, sellSide []priceVol) string {
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