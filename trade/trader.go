package trade

import (
	"math/rand"
	"sync/atomic"
)

type Market interface {
	Sell(offer int) (price int, sold bool)
	Buy(offer int) (price int, bought bool)
}

func atLeast(x, least int) int {
	if least > x {
		return least
	} else {
		return x
	}
}

func New(market Market) *Trader {
	t := &Trader{
		confidence: 50,
		assets:     rand.Int31n(10),
	}
	go func() {
		expect := 1000 + rand.Intn(100)
		sell := t.Assets() >= 5
		hist := []int{}
		for {
			if n := len(hist); n > 2 {
				boom := hist[n-1] > hist[n-2] && hist[n-2] > hist[n-3]
				sell = boom || t.Assets() > 5
			}
			if sell {
				if price, sold := market.Sell(expect); sold {
					hist = append(hist, price)
					t.moreConfident()
					expect += (price - expect) * t.Confidence() / 100
					t.decAssets()
				} else {
					t.lessConfident()
					margin := 100 * t.Confidence() / 100
					expect = atLeast(expect-margin, 1)
				}
			} else {
				if price, bought := market.Buy(expect); bought {
					hist = append(hist, price)
					t.moreConfident()
					expect += (price - expect) * t.Confidence() / 100
					t.incAssets()
				} else {
					t.lessConfident()
					margin := 100 * t.Confidence() / 100
					expect += margin
				}
			}
		}
	}()
	return t
}

type Trader struct {
	confidence int32
	assets     int32
}

func (t *Trader) Confidence() int {
	return int(atomic.LoadInt32(&t.confidence))
}

func (t *Trader) Assets() int {
	return int(atomic.LoadInt32(&t.assets))
}

func (t *Trader) incAssets() {
	atomic.AddInt32(&t.assets, 1)
}

func (t *Trader) decAssets() {
	atomic.AddInt32(&t.assets, -1)
}

func (t *Trader) moreConfident() {
	v := atomic.AddInt32(&t.confidence, 1)
	if v > 100 {
		atomic.CompareAndSwapInt32(&t.confidence, v, 100)
	}
}

func (t *Trader) lessConfident() {
	v := atomic.AddInt32(&t.confidence, -1)
	if v < 0 {
		atomic.CompareAndSwapInt32(&t.confidence, v, 0)
	}
}
