package filters

import (
	"net"

	"github.com/jlorgal/odor/odor"
)

// ParentalControl filter.
type ParentalControl struct {
	blacklist []*net.IPNet
}

// NewParentalControl creates a ParentalControl filter
func NewParentalControl(config *odor.Config) (*ParentalControl, error) {
	blacklist, err := odor.GetBlacklist("parental", config)
	return &ParentalControl{blacklist: blacklist}, err
}

// Request filters ingress packets.
func (p *ParentalControl) Request(context *odor.Context) odor.FilterAction {
	if context.Profile == nil || !context.Profile.ParentalControl {
		return odor.Accept
	}
	if ipv4 := odor.GetIPv4Layer(context.Packet); ipv4 != nil {
		if odor.IsBlacklistedIP(p.blacklist, ipv4.DstIP) {
			return odor.Drop
		}
	}
	return odor.Accept
}

// Response filters egress packets.
func (p *ParentalControl) Response(context *odor.Context) odor.FilterAction {
	return odor.Accept
}
