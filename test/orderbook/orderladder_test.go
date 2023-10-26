package orderbook

import (
	"reflect"
	"testing"

	"github.com/mbcolwell/golang/internal/orderbook"
)

func TestBisectDepth(t *testing.T) {
	pv := orderbook.PriceVol{Price: 500, Volume: 0}
	allSmaller := []orderbook.PriceVol{{Price: 100, Volume: 0}, {Price: 200, Volume: 0}, {Price: 300, Volume: 0}}
	allLarger := []orderbook.PriceVol{{Price: 501, Volume: 0}, {Price: 502, Volume: 0}, {Price: 503, Volume: 0}}
	middle := []orderbook.PriceVol{{Price: 100, Volume: 0}, {Price: 501, Volume: 0}, {Price: 1000, Volume: 0}}
	exact := []orderbook.PriceVol{{Price: 100, Volume: 0}, {Price: 200, Volume: 0}, {Price: 500, Volume: 0}}

	if orderbook.BisectDepth(pv.Price, allSmaller) != len(allSmaller) {
		t.Fatalf("Bisect depth function failed for all existing array elements smaller than new element")
	}
	if orderbook.BisectDepth(pv.Price, allLarger) != 0 {
		t.Fatalf("Bisect depth function failed for all existing array elements larger than new element")
	}
	if orderbook.BisectDepth(pv.Price, middle) != 1 {
		t.Fatalf("Bisect depth function failed for middle insertion of new element")
	}
	if orderbook.BisectDepth(pv.Price, exact) != 2 {
		t.Fatalf("Bisect depth function failed for matching existing array element")
	}
}

func TestOrderLadder(t *testing.T) {
	orders := []orderbook.Message{
		{
			Header: orderbook.Header{MsgSize: 0, SeqNo: 1},
			Order:  orderbook.Order{MsgType: byte('A'), Symbol: [3]byte([]byte("ABC")), OrderId: 0, Side: byte('B')},
			Size:   1,
			Price:  100,
		},
		{
			Header: orderbook.Header{MsgSize: 0, SeqNo: 2},
			Order:  orderbook.Order{MsgType: byte('A'), Symbol: [3]byte([]byte("ABC")), OrderId: 1, Side: byte('B')},
			Size:   50,
			Price:  101,
		},
		{
			Header: orderbook.Header{MsgSize: 0, SeqNo: 3},
			Order:  orderbook.Order{MsgType: byte('A'), Symbol: [3]byte([]byte("ABC")), OrderId: 2, Side: byte('B')},
			Size:   25,
			Price:  100,
		},
	}

	testCases := []struct {
		name           string
		runOrders      []int
		expectedLadder orderbook.Ladder
		expectedIdxs   []int
	}{
		{
			name:      "simple insert",
			runOrders: []int{0},
			expectedLadder: orderbook.Ladder{
				Depth:  &[]orderbook.PriceVol{{Price: 100, Volume: 1}},
				Orders: map[uint64]orderbook.PriceVol{0: {Price: 100, Volume: 1}}},
			expectedIdxs: []int{0},
		},
		{
			name:      "same price insert",
			runOrders: []int{0, 1, 2},
			expectedLadder: orderbook.Ladder{
				Depth: &[]orderbook.PriceVol{{Price: 100, Volume: 26}, {Price: 101, Volume: 50}},
				Orders: map[uint64]orderbook.PriceVol{
					0: {Price: 100, Volume: 1}, 1: {Price: 101, Volume: 50}, 2: {Price: 100, Volume: 25},
				}},
			expectedIdxs: []int{0, 1, 0},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var l orderbook.Ladder
			l.Depth = &[]orderbook.PriceVol{}
			l.Orders = map[uint64]orderbook.PriceVol{}

			idxResponses := []int{}
			var idxResponse int

			for _, orderIdx := range tc.runOrders {
				order := orders[orderIdx]
				switch order.Order.MsgType {
				case byte('A'):
					idxResponse = l.AddOrder(order.Order.OrderId, order.Price, order.Size)
					break
				case byte('U'):
					idxResponse = l.UpdateOrder(
						order.Order.OrderId, order.Price, order.Size, order.Order.Side,
					)
					break
				case byte('D'):
					idxResponse = l.DeleteOrder(order.Order.OrderId)
					break
				case byte('E'):
					idxResponse = l.ExecuteOrder(order.Order.OrderId, order.Size)
					break
				}
				idxResponses = append(idxResponses, idxResponse)
			}

			if !reflect.DeepEqual(*tc.expectedLadder.Depth, *l.Depth) {
				t.Fatalf(
					"Test failed on ladder comparison\nExpected: %#v\nOutput: %#v", *tc.expectedLadder.Depth, *l.Depth,
				)
			}
			if !reflect.DeepEqual(tc.expectedLadder.Orders, l.Orders) {
				t.Fatalf(
					"Test failed on ladder comparison\nExpected: %#v\nOutput: %#v", tc.expectedLadder.Orders, l.Orders,
				)
			}
			if !reflect.DeepEqual(tc.expectedIdxs, idxResponses) {
				t.Fatalf(
					"Test failed on index response comparison\nExpected: %v\nOutput: %v", tc.expectedIdxs, idxResponses,
				)
			}
		})
	}
}
