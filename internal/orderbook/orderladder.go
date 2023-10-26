package orderbook

type PriceVol struct {
	Price  int32
	Volume uint64
}

type Ladder struct {
	Depth  *[]PriceVol
	Orders map[uint64]PriceVol
}

func BisectDepth(price int32, pv []PriceVol) int {
	lo := 0
	hi := len(pv)
	for lo < hi {
		mid := (lo + hi) / 2
		if pv[mid].Price < price {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func (ladder *Ladder) AddOrder(id uint64, price int32, volume uint64) int {
	(*ladder).Orders[id] = PriceVol{Price: price, Volume: volume}
	if len((*(*ladder).Depth)) == 0 {
		(*(*ladder).Depth) = append((*(*ladder).Depth)[:], PriceVol{Price: price, Volume: volume})
		return 0
	}

	idx := BisectDepth(price, (*(*ladder).Depth))
	if idx == len((*(*ladder).Depth)) {
		(*(*ladder).Depth) = append((*(*ladder).Depth), PriceVol{Price: price, Volume: volume})
	} else if (*(*ladder).Depth)[idx].Price == price {
		(*(*ladder).Depth)[idx].Volume += volume
	} else if (*(*ladder).Depth)[idx].Price > price {
		(*(*ladder).Depth) = append((*(*ladder).Depth)[:idx+1], (*(*ladder).Depth)[idx:]...)
		(*(*ladder).Depth)[idx] = PriceVol{Price: price, Volume: volume}
	}
	return idx
}

func (ladder *Ladder) UpdateOrder(id uint64, price int32, volume uint64, side byte) int {
	var idx int
	if _, ok := (*ladder).Orders[id]; ok {
		idx1 := ladder.DeleteOrder(id)
		idx2 := ladder.AddOrder(id, price, volume)
		if side == byte('B') {
			return max(idx1, idx2)
		} else {
			return min(idx1, idx2)
		}

	} else {
		idx = ladder.AddOrder(id, price, volume) // Order was previously processed but had zero size
	}
	return idx
}

func (ladder *Ladder) DeleteOrder(id uint64) int {
	if existingOrder, ok := (*ladder).Orders[id]; ok {
		idx := BisectDepth(existingOrder.Price, (*(*ladder).Depth))
		if existingOrder.Volume == (*(*ladder).Depth)[idx].Volume {
			if idx == len((*(*ladder).Depth))-1 {
				(*(*ladder).Depth) = (*(*ladder).Depth)[:idx]
			} else {
				(*(*ladder).Depth) = append((*(*ladder).Depth)[:idx], (*(*ladder).Depth)[idx+1:]...)
			}
		} else {
			(*(*ladder).Depth)[idx].Volume -= existingOrder.Volume
		}
		delete((*ladder).Orders, id)
		return idx
	}
	return -1 // Order was never added before deletion (added with 0 size)
}

func (ladder *Ladder) ExecuteOrder(id uint64, size uint64) int {
	order, _ := (*ladder).Orders[id]
	if order.Volume == size {
		return ladder.DeleteOrder(id)
	} else {
		idx := BisectDepth(order.Price, (*(*ladder).Depth))
		(*(*ladder).Depth)[idx].Volume -= size
		(*ladder).Orders[id] = PriceVol{Price: (*ladder).Orders[id].Price, Volume: ((*ladder).Orders[id].Volume - size)}
		return idx
	}
}
