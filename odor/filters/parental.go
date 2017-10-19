package filters

import (
	"net"

	"github.com/jlorgal/odor/odor"
)

// ParentalControl filter.
type ParentalControl struct {
}

// NewParentalControl creates a ParentalControl filter
func NewParentalControl() *ParentalControl {
	return &ParentalControl{}
}

// Request filters ingress packets.
func (p *ParentalControl) Request(context *odor.Context) odor.FilterAction {
	if ipv4 := odor.GetIPv4Layer(context.Packet); ipv4 != nil {
		// TODO: Tomorrow we introduce machine learning :P
		if ipv4.DstIP.Equal(net.ParseIP("176.34.179.218")) {
			return odor.Drop
		}
	}
	return odor.Accept
}

// Response filters egress packets.
func (p *ParentalControl) Response(context *odor.Context) odor.FilterAction {
	return odor.Accept
}
