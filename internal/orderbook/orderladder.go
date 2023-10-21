package orderbook

type priceVol struct {
	Price  int32
	Volume uint64
}

type Ladder struct {
	Depth  []priceVol
	Orders map[uint64]priceVol
}

func bisectDepth(price int32, pv []priceVol) int {
	lo := 0
	hi := len(pv)
	for lo < hi {
		mid := (lo + hi) / 2
		if price < pv[mid].Price {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func addOrder(id uint64, price int32, volume uint64, ladder *Ladder) int {
	(*ladder).Orders[id] = priceVol{Price: price, Volume: volume}
	idx := bisectDepth(price, (*ladder).Depth)
	if (*ladder).Depth[idx].Price == price {
		(*ladder).Depth[idx].Volume += volume
	} else if (*ladder).Depth[idx].Price > price {
		newRecord := priceVol{Price: price, Volume: volume}
		if idx == len((*ladder).Depth) {
			(*ladder).Depth = append((*ladder).Depth, newRecord)
		} else {
			(*ladder).Depth = append((*ladder).Depth[:idx+1], (*ladder).Depth[idx:]...)
			(*ladder).Depth[idx] = newRecord
		}
	}
	return idx
}

func updateOrder(id uint64, price int32, volume uint64, ladder *Ladder) int {
	var idx int
	if existingOrder, ok := (*ladder).Orders[id]; ok {
		idx = bisectDepth(price, (*ladder).Depth)
		(*ladder).Depth[idx].Volume += (volume - existingOrder.Volume)
	} else {
		idx = addOrder(id, price, volume, ladder) // Order was previously processed but had zero size
	}
	return idx
}

func deleteOrder(id uint64, ladder *Ladder) int {
	if existingOrder, ok := (*ladder).Orders[id]; ok {
		idx := bisectDepth(existingOrder.Price, (*ladder).Depth)
		if existingOrder.Volume == (*ladder).Depth[idx].Volume {
			if idx == len((*ladder).Depth) {
				(*ladder).Depth = (*ladder).Depth[:idx]
			} else {
				(*ladder).Depth = append((*ladder).Depth[:idx], (*ladder).Depth[:idx+1]...)
			}
		} else {
			(*ladder).Depth[idx].Volume -= existingOrder.Volume
		}
		delete((*ladder).Orders, id)
		return idx
	}
	return -1 // Order was never added before deletion (added with 0 size)
}

func executeOrder(id uint64, size uint64, ladder *Ladder) int {
	order, _ := (*ladder).Orders[id]
	if order.Volume == size {
		return deleteOrder(id, ladder)
	} else {
		idx := bisectDepth(order.Price, (*ladder).Depth)
		(*ladder).Depth[idx].Volume -= size
		return idx
	}
}
