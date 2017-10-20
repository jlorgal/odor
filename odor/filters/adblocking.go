package filters

import (
	"net"

	"github.com/jlorgal/odor/odor"
)

// AdBlocking filter.
type AdBlocking struct {
	blacklist []*net.IPNet
}

// NewAdBlocking creates a Malware filter
func NewAdBlocking(config *odor.Config) (*AdBlocking, error) {
	blacklist, err := odor.GetBlacklist("adBlocking", config)
	return &AdBlocking{blacklist: blacklist}, err
}

// Request filters ingress packets.
func (a *AdBlocking) Request(context *odor.Context) odor.FilterAction {
	if context.Profile == nil || !context.Profile.AdBlocking {
		return odor.Accept
	}
	if ipv4 := odor.GetIPv4Layer(context.Packet); ipv4 != nil {
		if odor.IsBlacklistedIP(a.blacklist, ipv4.DstIP) {
			return odor.Drop
		}
	}
	return odor.Accept
}

// Response filters egress packets.
func (a *AdBlocking) Response(context *odor.Context) odor.FilterAction {
	return odor.Accept
}
