package binomv2postback

type PostbackLevel uint8

const (
	PB_LVL_ALL      PostbackLevel = iota // postback to tracker and traffic source
	PB_LVL_NO_TS                         //
	PB_LVL_DISABLED                      // disable postback (dryRun?)
)

func (lvl PostbackLevel) CanPostback() bool {
	return lvl < PB_LVL_DISABLED
}

func (lvl PostbackLevel) String() string {
	switch lvl {
	case PB_LVL_NO_TS:
		return "traffic_source_postback_disabled"
	case PB_LVL_DISABLED:
		return "postback_disabled"
	default:
		return "postback_all"
	}
}
