package app

type (
	Document struct {
		// ProductID.
		Id int
		// URL
		URL string
		// Subject.
		Title string
		// TradeCount.
		TotalSales int
		// StarSTR.
		RatingProduct float64
		// CountFeedback.
		TotalComment int
		// ImagePathList
		Images []string

		// Discount.
		Discount float64
		// Max.
		Max Price
		// Min.
		Min Price

		// Store.
		Shop Shop
	}

	Price struct {
		// Currency.
		Currency string
		// Value.
		Cost float64
	}

	Shop struct {
		// StoreID.
		ID int
		// StoreName.
		Name string
		// FollowingNumber.
		Followers int
		// PositiveRateSTR.
		PositiveRate float64
	}
)

func (d *Document) Rating() float64 {
	return reviewsRate(d.TotalComment) *
		(0.1*favoritesRate(9999) +
			0.2*ordersRate(d.TotalSales) +
			0.5*costRate(d.Min.Cost) +
			0.15*sellerRate(d.Shop.PositiveRate) +
			0.05*subscribersRate(d.Shop.Followers))
}

func sellerRate(rateSeller float64) (rate float64) {
	switch {
	case rateSeller <= 90:
		rate = 0.0
	case rateSeller > 90 && rateSeller < 95:
		rate = 0.5
	case rateSeller >= 95:
		rate = 1.0
	}
	return rate
}

func reviewsRate(reviewsCount int) (rate float64) {
	switch {
	case reviewsCount <= 10:
		rate = 0.0
	case reviewsCount > 10 && reviewsCount < 100:
		rate = 0.3
	case reviewsCount >= 100:
		rate = 1.0
	}
	return rate
}

func favoritesRate(favoritesCount int) (rate float64) {
	switch {
	case favoritesCount >= 9999:
		rate = 1.0
	case favoritesCount < 9999 && favoritesCount >= 5000:
		rate = 0.7
	case favoritesCount < 5000 && favoritesCount >= 1000:
		rate = 0.5
	case favoritesCount < 1000 && favoritesCount >= 100:
		rate = 0.2
	case favoritesCount < 100:
		rate = 0.1
	}
	return rate
}

func ordersRate(ordersCount int) (rate float64) {
	switch {
	case ordersCount >= 1000:
		rate = 1.0
	case ordersCount < 1000 && ordersCount >= 300:
		rate = 0.7
	case ordersCount < 300 && ordersCount >= 100:
		rate = 0.4
	case ordersCount < 100:
		rate = 0.1
	}
	return rate
}

func costRate(cost float64) (rate float64) {
	switch {
	case cost >= 1.0 && cost <= 5.0:
		rate = 1.0
	case cost <= 10.0 && cost > 5.0:
		rate = 0.7
	case cost <= 100.0 && cost > 10.0:
		rate = 0.5
	case cost < 1.0 || cost > 100.0:
		rate = 0.35
	}
	return rate
}

func subscribersRate(subscribersCount int) (rate float64) {
	switch {
	case subscribersCount >= 100000:
		rate = 1.0
	case subscribersCount < 100000 && subscribersCount >= 10000:
		rate = 0.7
	case subscribersCount < 10000 && subscribersCount >= 1000:
		rate = 0.5
	case subscribersCount < 1000 && subscribersCount >= 100:
		rate = 0.2
	case subscribersCount < 100:
		rate = 0.1
	}
	return rate
}
