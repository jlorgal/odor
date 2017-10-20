package filters

import (
	"github.com/jlorgal/odor/odor"
	"github.com/jlorgal/odor/odor/profile"
)

// IdentifyUser filter.
type IdentifyUser struct {
}

// IdentifyUser creates a IdentifyUser filter
func NewIdentifyUser() *IdentifyUser {
	return &IdentifyUser{}
}

// Request filters ingress packets.
func (p *IdentifyUser) Request(context *odor.Context) odor.FilterAction {

	if ipv4 := odor.GetIPv4Layer(context.Packet); ipv4 != nil {
		srcIP := ipv4.SrcIP
		if radiusPacket, err := profile.GetRadiusPacket(srcIP.String()); err == nil {
			context.Profile.MSISDN = radiusPacket.MSISDN
		}
	}
	return odor.Accept
}

// Response filters egress packets.
func (p *IdentifyUser) Response(context *odor.Context) odor.FilterAction {
	return odor.Accept
}
