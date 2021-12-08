package stock

func firstMatch(buyers, sellers []*transaction) {
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
