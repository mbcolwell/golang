package orderbook

import (
	"reflect"
	"strings"
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
		orderbook.CreateOrder(0, 'A', "ABC", 0, 'B', 1, 100),
		orderbook.CreateOrder(1, 'A', "ABC", 1, 'B', 50, 101),
		orderbook.CreateOrder(2, 'A', "ABC", 2, 'B', 25, 100),
		orderbook.CreateOrder(3, 'U', "ABC", 0, 'B', 30, 102),
		orderbook.CreateOrder(4, 'D', "ABC", 1, 'B', 0, 0),
		orderbook.CreateOrder(5, 'D', "ABC", 2, 'B', 0, 0),
		orderbook.CreateOrder(6, 'E', "ABC", 0, 'B', 10, 0),
		orderbook.CreateOrder(7, 'E', "ABC", 0, 'B', 30, 0),
		orderbook.CreateOrder(8, 'A', "ABC", 3, 'S', 30, 99),
		orderbook.CreateOrder(9, 'A', "ABC", 4, 'S', 30, 101),
		orderbook.CreateOrder(10, 'A', "ABC", 5, 'S', 30, 100),
		orderbook.CreateOrder(11, 'U', "ABC", 5, 'S', 30, 102),
	}

	testCases := []struct {
		name         string
		runOrders    []int
		Depth        []orderbook.PriceVol
		Orders       map[uint64]orderbook.PriceVol
		expectedIdxs []int
	}{
		{
			name:         "simple insert",
			runOrders:    []int{0},
			Depth:        []orderbook.PriceVol{{100, 1}},
			Orders:       map[uint64]orderbook.PriceVol{0: {100, 1}},
			expectedIdxs: []int{0},
		},
		{
			name:         "same price insert",
			runOrders:    []int{0, 1, 2},
			Depth:        []orderbook.PriceVol{{100, 26}, {101, 50}},
			Orders:       map[uint64]orderbook.PriceVol{0: {100, 1}, 1: {101, 50}, 2: {100, 25}},
			expectedIdxs: []int{0, 1, 0},
		},
		{
			name:         "simple update",
			runOrders:    []int{0, 1, 2, 3},
			Depth:        []orderbook.PriceVol{{100, 25}, {101, 50}, {102, 30}},
			Orders:       map[uint64]orderbook.PriceVol{0: {102, 30}, 1: {101, 50}, 2: {100, 25}},
			expectedIdxs: []int{0, 1, 0, 2},
		},
		{
			name:         "update clears price",
			runOrders:    []int{0, 1, 3},
			Depth:        []orderbook.PriceVol{{101, 50}, {102, 30}},
			Orders:       map[uint64]orderbook.PriceVol{0: {102, 30}, 1: {101, 50}},
			expectedIdxs: []int{0, 1, 1},
		},
		{
			name:         "update sell side",
			runOrders:    []int{8, 9, 10, 11},
			Depth:        []orderbook.PriceVol{{99, 30}, {101, 30}, {102, 30}},
			Orders:       map[uint64]orderbook.PriceVol{3: {99, 30}, 4: {101, 30}, 5: {102, 30}},
			expectedIdxs: []int{0, 1, 1, 1},
		},
		{
			name:         "simple delete",
			runOrders:    []int{1, 4},
			Depth:        []orderbook.PriceVol{},
			Orders:       map[uint64]orderbook.PriceVol{},
			expectedIdxs: []int{0, 0},
		},
		{
			name:         "delete with remaining at price",
			runOrders:    []int{0, 2, 5},
			Depth:        []orderbook.PriceVol{{100, 1}},
			Orders:       map[uint64]orderbook.PriceVol{0: {100, 1}},
			expectedIdxs: []int{0, 0, 0},
		},
		{
			name:         "partial execution",
			runOrders:    []int{0, 3, 6},
			Depth:        []orderbook.PriceVol{{102, 20}},
			Orders:       map[uint64]orderbook.PriceVol{0: {102, 20}},
			expectedIdxs: []int{0, 0, 0},
		},
		{
			name:         "full execution",
			runOrders:    []int{0, 3, 7},
			Depth:        []orderbook.PriceVol{},
			Orders:       map[uint64]orderbook.PriceVol{},
			expectedIdxs: []int{0, 0, 0},
		},
		{
			name:         "multi-step execution",
			runOrders:    []int{0, 3, 6, 6, 6},
			Depth:        []orderbook.PriceVol{},
			Orders:       map[uint64]orderbook.PriceVol{},
			expectedIdxs: []int{0, 0, 0, 0, 0},
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
				case byte('U'):
					idxResponse = l.UpdateOrder(
						order.Order.OrderId, order.Price, order.Size, order.Order.Side,
					)
				case byte('D'):
					idxResponse = l.DeleteOrder(order.Order.OrderId)
				case byte('E'):
					idxResponse = l.ExecuteOrder(order.Order.OrderId, order.Size)
				}
				idxResponses = append(idxResponses, idxResponse)
			}

			if !reflect.DeepEqual(tc.Depth, *l.Depth) {
				t.Fatalf(
					"Test failed on ladder comparison\nExpected: %#v\nOutput: %#v", tc.Depth, *l.Depth,
				)
			}
			if !reflect.DeepEqual(tc.Orders, l.Orders) {
				t.Fatalf(
					"Test failed on ladder comparison\nExpected: %#v\nOutput: %#v", tc.Orders, l.Orders,
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

func TestOrderBook(t *testing.T) {
	orders := []orderbook.Message{
		orderbook.CreateOrder(0, 'A', "ABC", 0, 'B', 50, 100),
		orderbook.CreateOrder(1, 'A', "ABC", 1, 'B', 50, 101),
		orderbook.CreateOrder(2, 'A', "ABC", 2, 'B', 50, 102),
		orderbook.CreateOrder(3, 'A', "DEF", 3, 'S', 50, 100),
		orderbook.CreateOrder(4, 'A', "DEF", 4, 'S', 50, 100),
		orderbook.CreateOrder(5, 'A', "DEF", 5, 'S', 0, 100),
		orderbook.CreateOrder(6, 'D', "DEF", 4, 'S', 0, 0),
		orderbook.CreateOrder(7, 'U', "DEF", 3, 'S', 30, 100),
		orderbook.CreateOrder(8, 'E', "DEF", 3, 'S', 1, 0),
		orderbook.CreateOrder(9, 'E', "DEF", 3, 'S', 28, 0),
		orderbook.CreateOrder(10, 'E', "DEF", 3, 'S', 1, 0),
		orderbook.CreateOrder(11, 'A', "DEF", 6, 'B', 0, 100),
		orderbook.CreateOrder(12, 'U', "DEF", 6, 'B', 50, 100),
		orderbook.CreateOrder(13, 'A', "ABC", 7, 'B', 50, 103),
		orderbook.CreateOrder(14, 'A', "ABC", 8, 'B', 50, 104),
		orderbook.CreateOrder(15, 'A', "ABC", 9, 'B', 50, 105),
		orderbook.CreateOrder(16, 'U', "ABC", 0, 'B', 100, 100),
		orderbook.CreateOrder(17, 'A', "ABC", 10, 'B', 50, 99),
		orderbook.CreateOrder(18, 'E', "ABC", 9, 'B', 50, 0),
		orderbook.CreateOrder(19, 'D', "ABC", 8, 'B', 0, 0),
		orderbook.CreateOrder(20, 'A', "GHI", 11, 'S', 50, 500),
		orderbook.CreateOrder(21, 'A', "GHI", 12, 'S', 50, 501),
		orderbook.CreateOrder(22, 'A', "GHI", 13, 'S', 50, 502),
		orderbook.CreateOrder(23, 'A', "GHI", 14, 'S', 50, 503),
		orderbook.CreateOrder(24, 'A', "GHI", 15, 'S', 50, 504),
		orderbook.CreateOrder(25, 'A', "GHI", 16, 'S', 50, 505),
		orderbook.CreateOrder(26, 'U', "GHI", 11, 'S', 100, 500),
		orderbook.CreateOrder(27, 'A', "GHI", 17, 'S', 50, 499),
		orderbook.CreateOrder(28, 'A', "GHI", 18, 'S', 50, 499),
		orderbook.CreateOrder(29, 'E', "GHI", 17, 'S', 50, 0),
		orderbook.CreateOrder(30, 'D', "GHI", 13, 'S', 0, 0),
		orderbook.CreateOrder(31, 'E', "GHI", 18, 'S', 50, 0),
		orderbook.CreateOrder(32, 'D', "GHI", 14, 'S', 0, 0),
		orderbook.CreateOrder(33, 'A', "GHI", 19, 'S', 0, 499),
		orderbook.CreateOrder(34, 'D', "GHI", 19, 'S', 0, 0),
		orderbook.CreateOrder(35, 'A', "GHI", 20, 'B', 2000, 350),
		orderbook.CreateOrder(36, 'A', "GHI", 21, 'B', 2, 351),
	}

	expected_lines := []string{
		"0, ABC, [(100, 50)], []",
		"1, ABC, [(101, 50), (100, 50)], []",
		"2, ABC, [(102, 50), (101, 50), (100, 50)], []",
		"3, DEF, [], [(100, 50)]",
		"4, DEF, [], [(100, 100)]",
		"6, DEF, [], [(100, 50)]",
		"7, DEF, [], [(100, 30)]",
		"8, DEF, [], [(100, 29)]",
		"9, DEF, [], [(100, 1)]",
		"10, DEF, [], []",
		"12, DEF, [(100, 50)], []",
		"13, ABC, [(103, 50), (102, 50), (101, 50), (100, 50)], []",
		"14, ABC, [(104, 50), (103, 50), (102, 50), (101, 50), (100, 50)], []",
		"15, ABC, [(105, 50), (104, 50), (103, 50), (102, 50), (101, 50)], []",
		"18, ABC, [(104, 50), (103, 50), (102, 50), (101, 50), (100, 100)], []",
		"19, ABC, [(103, 50), (102, 50), (101, 50), (100, 100), (99, 50)], []",
		"20, GHI, [], [(500, 50)]",
		"21, GHI, [], [(500, 50), (501, 50)]",
		"22, GHI, [], [(500, 50), (501, 50), (502, 50)]",
		"23, GHI, [], [(500, 50), (501, 50), (502, 50), (503, 50)]",
		"24, GHI, [], [(500, 50), (501, 50), (502, 50), (503, 50), (504, 50)]",
		"26, GHI, [], [(500, 100), (501, 50), (502, 50), (503, 50), (504, 50)]",
		"27, GHI, [], [(499, 50), (500, 100), (501, 50), (502, 50), (503, 50)]",
		"28, GHI, [], [(499, 100), (500, 100), (501, 50), (502, 50), (503, 50)]",
		"29, GHI, [], [(499, 50), (500, 100), (501, 50), (502, 50), (503, 50)]",
		"30, GHI, [], [(499, 50), (500, 100), (501, 50), (503, 50), (504, 50)]",
		"31, GHI, [], [(500, 100), (501, 50), (503, 50), (504, 50), (505, 50)]",
		"32, GHI, [], [(500, 100), (501, 50), (504, 50), (505, 50)]",
		"35, GHI, [(350, 2000)], [(500, 100), (501, 50), (504, 50), (505, 50)]",
		"36, GHI, [(351, 2), (350, 2000)], [(500, 100), (501, 50), (504, 50), (505, 50)]",
		"",
	}
	expected := strings.Join(expected_lines, "\n")

	var output string
	var ticker string
	book := map[string]orderbook.Ladder{}
	n := 5

	for _, order := range orders {
		if orderbook.ProcessMessage(n, order, book) {
			ticker = string(order.Order.Symbol[:])
			output += orderbook.FormatLadder(
				n, ticker, order.Header.SeqNo, *book[ticker+"B"].Depth, *book[ticker+"S"].Depth,
			)
			output += "\n"
		}
	}
	if output != expected{
		t.Fatalf("Expected:\n%s\n\nOutput:\n%s", expected, output)
	}
}
