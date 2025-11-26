package binomv2postback

type Postback struct {
	ID     string   `json:"cnv_id"`
	Events Events   `json:"events,omitempty"`
	Payout *float64 `json:"payout,omitempty"`
	Status *string  `json:"cnv_status,omitempty"`
}
