package odor

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
)

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

// GetBlacklist generates a blacklist of CIDRs for a filter from the configuration.
func GetBlacklist(filter string, config *Config) ([]*net.IPNet, error) {
	blacklist := []*net.IPNet{}
	if config.Filters[filter] == nil {
		return blacklist, nil
	}
	for _, ipnet := range config.Filters[filter] {
		_, net, err := net.ParseCIDR(ipnet)
		if err != nil {
			return blacklist, fmt.Errorf("Invalid blacklist element: %s. %s", ipnet, err)
		}
		blacklist = append(blacklist, net)
	}
	return blacklist, nil
}

// IsBlacklistedIP checks if an IP address is blacklisted.
func IsBlacklistedIP(blacklist []*net.IPNet, ip net.IP) bool {
	for _, ipnet := range blacklist {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}
