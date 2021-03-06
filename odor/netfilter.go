package odor

import (
	"syscall"

	"github.com/chifflier/nfqueue-go/nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var netFilterStatic *NetFilter

// PacketHandler interface that implements the callback to handle a packet.
type PacketHandler interface {
	HandlePacket(context *Context) FilterAction
}

// NetFilter provides some operations over a NetFilter queue.
type NetFilter struct {
	queue   *nfqueue.Queue
	handler PacketHandler
}

// NewNetFilter creates a NetFilter object.
func NewNetFilter(handler PacketHandler) *NetFilter {
	queue := new(nfqueue.Queue)
	netFilter := &NetFilter{
		queue:   queue,
		handler: handler,
	}
	netFilterStatic = netFilter
	queue.SetCallback(staticCallback)
	return netFilter
}

func staticCallback(payload *nfqueue.Payload) int {
	return netFilterStatic.Callback(payload)
}

// Start the NetFilter queue.
func (n *NetFilter) Start(queueNum int) {
	n.queue.Init()
	n.queue.Unbind(syscall.SOCK_PACKET)
	n.queue.Bind(syscall.SOCK_PACKET)
	n.queue.CreateQueue(queueNum)
	n.queue.Loop()
}

// Stop the NetFilter queue.
func (n *NetFilter) Stop() {
	n.queue.DestroyQueue()
	n.queue.Close()
}

// Callback to handle a packet from NetFilter queue.
func (n NetFilter) Callback(payload *nfqueue.Payload) int {
	// Decode a packet
	packet := gopacket.NewPacket(payload.Data, layers.LayerTypeIPv4, gopacket.Default)
	context := &Context{
		PacketInitial: packet,
		Packet:        packet,
	}
	action := n.handler.HandlePacket(context)
	switch action {
	case Accept:
		payload.SetVerdict(nfqueue.NF_ACCEPT)
	case Drop:
		payload.SetVerdict(nfqueue.NF_DROP)
	}
	return 0
}

// GetIPv4Layer returns the IPv4 layer of the packet.
func GetIPv4Layer(packet gopacket.Packet) *layers.IPv4 {
	if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
		return ipv4Layer.(*layers.IPv4)
	}
	return nil
}
