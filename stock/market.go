package stock

import (
	"sync"
	"time"
)

func New(closeInterval time.Duration) *Service {
	s := &Service{stats: &Stats{}}
	go func() {
		ticker := time.NewTicker(closeInterval)
		for range ticker.C {
			s.mutex.Lock()
			for i := 0; i < len(s.buyers); i++ {
				buyer := s.buyers[i]
				for _, seller := range s.sellers {
					if buyer.Offer >= seller.Offer && seller.Price == 0 {
						seller.Price = buyer.Offer
						buyer.Price = seller.Offer
						s.updateStats(seller.Offer)
						break
					}
				}
			}
			for _, b := range s.buyers {
				b.Barrier.Done()
			}
			for _, s := range s.sellers {
				s.Barrier.Done()
			}
			s.rounds++
			s.buyers = nil
			s.sellers = nil
			s.mutex.Unlock()
		}
		ticker.Stop()
	}()
	return s
}

type Stats struct {
	Open         int
	High         int
	Low          int
	Close        int
	Transactions int
}

type transaction struct {
	Offer   int
	Price   int
	Barrier sync.WaitGroup
}

type Service struct {
	buyers  []*transaction
	sellers []*transaction
	rounds  int
	stats   *Stats
	mutex   sync.Mutex
}

func (s *Service) Sell(price int) (int, bool) {
	tx := &transaction{Offer: price}
	tx.Barrier.Add(1)
	s.mutex.Lock()
	s.sellers = append(s.sellers, tx)
	s.mutex.Unlock()
	tx.Barrier.Wait()
	return tx.Price, tx.Price > 0
}

func (s *Service) Buy(price int) (int, bool) {
	tx := &transaction{Offer: price}
	tx.Barrier.Add(1)
	s.mutex.Lock()
	s.buyers = append(s.buyers, tx)
	s.mutex.Unlock()
	tx.Barrier.Wait()
	return tx.Price, tx.Price > 0
}

func (s *Service) Stat() (int, int, int, int, int) {
	s.mutex.Lock()
	cur := s.stats
	cur.Transactions /= s.rounds
	s.stats = &Stats{Open: cur.Close}
	s.rounds = 0
	s.mutex.Unlock()
	return cur.Open, cur.High, cur.Low, cur.Close, cur.Transactions
}

func (s *Service) updateStats(price int) {
	if s.stats.Open == 0 {
		s.stats.Open = price
	}
	if price > s.stats.High {
		s.stats.High = price
	}
	if price < s.stats.Low || s.stats.Low == 0 {
		s.stats.Low = price
	}
	s.stats.Close = price
	s.stats.Transactions++
}
