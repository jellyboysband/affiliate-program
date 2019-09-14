package app

type (
	Product struct {
		Id            string
		AliID         int
		Rating        float64
		URL           string
		Title         string
		TotalSales    int
		RatingProduct string
		TotalComment  int
		Discount      float64
		Max           Price
		Min           Price
		Shop          Shop
	}

	Price struct {
		Currency string
		Cost     float64
	}

	Shop struct {
		ID           int
		Name         string
		Followers    int
		PositiveRate string
	}

	ArgSaveProduct struct {
		AliID         int
		Rating        float64
		URL           string
		Title         string
		TotalSales    int
		RatingProduct string
		TotalComment  int
		Discount      float64
		Max           Price
		Min           Price
		Shop          Shop
	}
)
