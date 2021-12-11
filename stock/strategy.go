package stock

import (
	"math/rand"
	"sort"
)

var algs = map[string]func([]*Tx, []*Tx){
	"first":  firstMatch,
	"random": randomMatch,
	"sort":   sortMatch,
}

func Match(buyers, sellers []*Tx, alg string) {
	algs[alg](buyers, sellers)
}

func firstMatch(buyers, sellers []*Tx) {
	for i := 0; i < len(buyers); i++ {
		buyer := buyers[i]
		for _, seller := range sellers {
			if buyer.Offer >= seller.Offer && seller.Price == 0 {
				seller.Price = buyer.Offer
				buyer.Price = seller.Offer
				break
			}
		}
	}
}

func randomMatch(buyers, sellers []*Tx) {
	rand.Shuffle(len(sellers), func(i, j int) {
		sellers[i], sellers[j] = sellers[j], sellers[i]
	})
	for i := 0; i < len(buyers); i++ {
		buyer := buyers[i]
		for _, seller := range sellers {
			if buyer.Offer >= seller.Offer && seller.Price == 0 {
				seller.Price = buyer.Offer
				buyer.Price = seller.Offer
				break
			}
		}
	}
}

func sortMatch(buyers, sellers []*Tx) {
	sort.Slice(buyers, func(i, j int) bool {
		return buyers[i].Offer < buyers[j].Offer
	})
	sort.Slice(sellers, func(i, j int) bool {
		return sellers[i].Offer < sellers[j].Offer
	})
	for i := 0; i < len(buyers); i++ {
		buyer := buyers[i]
		for _, seller := range sellers {
			if buyer.Offer >= seller.Offer && seller.Price == 0 {
				seller.Price = buyer.Offer
				buyer.Price = seller.Offer
				break
			}
		}
	}
}
