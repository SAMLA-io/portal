package schema

type Subscription struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	CurrentPeriodEnd int64  `json:"current_period_end"`
	ProductID        string `json:"product_id"`
	PriceID          string `json:"price_id"`
}
