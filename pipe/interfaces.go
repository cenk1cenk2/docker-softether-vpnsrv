package pipe

type (
	DnsMasqConfigurationTemplate struct {
		TapInterface      string
		RangeStartAddress string
		RangeEndAddress   string
		RangeNetmask      string
		LeaseTime         string
		Gateway           string
		ForwardingZone    []string
		Options           map[string]string
	}
)
