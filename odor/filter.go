package odor

import "github.com/google/gopacket"

// Profile contains the user identity (msisdn) and settings.
type Profile struct {
	MSISDN          string
	AntiPhishing    bool
	AntiMalware     bool
	ParentalControl bool
	AdBlocking      bool
	Captive         bool
}

// Context to use by filters.
type Context struct {
	PacketInitial gopacket.Packet
	Packet        gopacket.Packet
	Profile       *Profile
}

// FilterAction type for enumeration of filter actions
type FilterAction int

const (
	// Accept the packet according to this filter
	Accept FilterAction = iota
	// Drop the packet according to this filter
	Drop
	// Modify the packet according to this filter
	Modify
)

// Filter interface to implement by step of the pipeline.
type Filter interface {
	Request(context *Context) FilterAction
	Response(context *Context) FilterAction
}
