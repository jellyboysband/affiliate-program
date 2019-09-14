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
		RatingProduct string
		// CountFeedback.
		TotalComment int

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
		PositiveRate string
	}
)
