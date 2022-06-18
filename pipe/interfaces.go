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
	}

	SoftEtherConfigurationTemplate struct {
		Interface  string
		DefaultHub string
	}
)
