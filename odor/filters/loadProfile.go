package filters

import (
	"github.com/jlorgal/odor/odor"
	"github.com/jlorgal/odor/odor/profile"
)

// LoadProfile filter.
type LoadProfile struct {
}

// NewLoadProfile creates a LoadProfile filter
func NewLoadProfile() *LoadProfile {
	return &LoadProfile{}
}

// Request filters ingress packets.
func (p *LoadProfile) Request(context *odor.Context) odor.FilterAction {

	if context.Profile != nil && context.Profile.MSISDN != nil {
		if profile, err := profile.GetUserProfile(context.Profile.MSISDN); err == nil {
			context.Profile = profile
		}
	}

	return odor.Accept
}

// Response filters egress packets.
func (p *LoadProfile) Response(context *odor.Context) odor.FilterAction {
	return odor.Accept
}
