package filters

import (
	"github.com/jlorgal/odor/odor"
	"github.com/jlorgal/odor/odor/profile"
	"github.com/jlorgal/odor/odor/svc"
)

// IdentifyUser filter.
type IdentifyUser struct {
}

// NewIdentifyUser creates a IdentifyUser filter
func NewIdentifyUser() *IdentifyUser {
	return &IdentifyUser{}
}

// Request filters ingress packets.
func (p *IdentifyUser) Request(context *odor.Context) odor.FilterAction {

	if ipv4 := odor.GetIPv4Layer(context.Packet); ipv4 != nil {
		srcIP := ipv4.SrcIP
		svc.NewLogger().Info("User IP address: %v", srcIP)
		if radiusPacket, err := profile.GetRadiusPacket(srcIP.String()); err == nil {
			context.Profile = &odor.Profile{MSISDN: radiusPacket.MSISDN}
		}
	}
	svc.NewLogger().Info("Identify user: %v", context.Profile)
	return odor.Accept
}

// Response filters egress packets.
func (p *IdentifyUser) Response(context *odor.Context) odor.FilterAction {
	return odor.Accept
}
