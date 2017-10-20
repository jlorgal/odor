package odor

import (
	"net"

	"github.com/jlorgal/odor/odor/svc"
)

// Profile contains the user identity (msisdn) and settings.
type Profile struct {
	MSISDN          string `json:"msisdn"`
	AntiPhishing    bool   `json:"antiPhising"`
	AntiMalware     bool   `json:"antiMalware"`
	ParentalControl bool   `json:"parentalControl"`
	AdBlocking      bool   `json:"adBlocking"`
	Captive         bool   `json:"captive"`
}

// RadiusPacket models a radius packet to map IP and MSISDN. IP is a string
type RadiusPacket struct {
	IP     string `json:"ip"`
	MSISDN string `json:"msisdn"`
}

// NetRadiusPacket models a radius packet to map IP and MSISDN. IP is a net.IP structure
type NetRadiusPacket struct {
	IP     net.IP `json:"ip"`
	MSISDN string `json:"msisdn"`
}

// RadiusPacket2NetRadiusPacket transforms a RadiusPacket into NetRadiusPacket
func RadiusPacket2NetRadiusPacket(packet RadiusPacket) (NetRadiusPacket, error) {
	netRP := NetRadiusPacket{}
	netRP.IP = net.ParseIP(packet.IP)
	if netRP.IP == nil {
		return netRP, svc.NewInvalidRequestError("Bad Request", "Bad Request")
	}
	netRP.MSISDN = packet.MSISDN
	return netRP, nil
}
