package parse

type (
	Page struct {
		TitleModule    TitleModule    `json:"titleModule"`
		PriceModule    PriceModule    `json:"priceModule"`
		StoreModule    StoreModule    `json:"storeModule"`
		RedirectModule RedirectModule `json:"redirectModule"`
		ActionModule   ActionModule   `json:"actionModule"`
		PageModule     PageModule     `json:"pageModule"`
		ImageModule    ImageModule    `json:"imageModule"`
	}

	TitleModule struct {
		Rating     FeedbackRating `json:"feedbackRating"`
		TradeCount int            `json:"tradeCount"`
		Subject    string         `json:"subject"`
	}

	FeedbackRating struct {
		StarSTR       string `json:"averageStar"`
		CountFeedback int    `json:"totalValidNum"`
	}

	StoreModule struct {
		FollowingNumber int    `json:"followingNumber"`
		PositiveRateSTR string `json:"positiveRate"`
		StoreName       string `json:"storeName"`
		StoreID         int    `json:"storeNum"`
	}

	PriceModule struct {
		Discount float64 `json:"discount"`
		Max      Amount  `json:"maxAmount"`
		Min      Amount  `json:"minActivityAmount"`
	}

	Amount struct {
		Currency string  `json:"currency"`
		Value    float64 `json:"value"`
	}

	PageModule struct {
		URL       string `json:"itemDetailUrl"`
		ProductID int    `json:"productId"`
	}

	ImageModule struct {
		ImagePathList []string `json:"imagePathList"`
	}
)

// Need for validate
type (
	RedirectModule struct {
		Code string `json:"code"`
	}
	ActionModule struct {
		ItemStatus         int `json:"itemStatus"`
		TotalAvailQuantity int `json:"totalAvailQuantity"`
	}
)
