package odor

import (
	"net"
	"testing"
)

var config = &Config{
	Filters: map[string][]string{
		"parental": []string{"216.58.192.0/19"},
	},
}

func TestGetBlacklist(t *testing.T) {
	blacklist, err := GetBlacklist("parental", config)
	if err != nil {
		t.Errorf("Error getting the parental blacklist. %s", err)
		return
	}
	if len(blacklist) != 1 {
		t.Errorf("Unexpected length for parental blacklist")
		return
	}
}

func TestIsBlacklistedIP(t *testing.T) {
	blacklist, _ := GetBlacklist("parental", config)
	tests := []struct {
		ip          string
		blacklisted bool
	}{
		{"216.58.192.3", true},
		{"215.58.192.3", false},
	}

	for _, test := range tests {
		blacklisted := IsBlacklistedIP(blacklist, net.ParseIP(test.ip))
		if blacklisted != test.blacklisted {
			t.Errorf("IP %s invalid blacklisted check. Expected: %v, got: %v", test.ip, test.blacklisted, blacklisted)
			return
		}
	}
}
