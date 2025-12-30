package binomv2postback

type Postback interface {
	ClickID() string
	Payout() string
	ConversionStatus() string
	ConversionStatus2() string
	Events() Events
	Params() []string
	String() string
}
