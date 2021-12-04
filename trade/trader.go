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
		mustSell := t.Assets() >= 5
		for {
			margin := atLeast(100*t.Confidence()/100, 1)
			if mustSell && t.Assets() > 0 {
				if price, sold := market.Sell(expect + margin); sold {
					t.moreConfident()
					expect = price
					t.decAssets()
				} else {
					t.lessConfident()
					expect = atLeast(expect-margin, 1)
					mustSell = t.Assets() >= 5
				}
			} else {
				if price, bought := market.Buy(expect - margin); bought {
					t.moreConfident()
					expect = price
					t.incAssets()
				} else {
					t.lessConfident()
					expect += margin
					mustSell = t.Assets() >= 5
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
