package app

type (
	Product struct {
		Id            string
		AliID         int
		Rating        float64
		Images        []string
		OurRating     float64
		URL           string
		Title         string
		TotalSales    int
		RatingProduct float64
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
		PositiveRate float64
	}

	ArgSaveProduct struct {
		Images        []string
		OurRating     float64
		AliID         int
		Rating        float64
		URL           string
		Title         string
		TotalSales    int
		RatingProduct float64
		TotalComment  int
		Discount      float64
		Max           Price
		Min           Price
		Shop          Shop
	}
)
