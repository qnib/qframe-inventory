package qframe_inventory

/******************** Inventory Request
 Sends a query for a key or an IP and provides a back-channel, so that the requesting partner can block on the request
 until it arrives - honouring a timeout...
*/

import (
	"time"
)

type Interface interface {
	Equal(interface{}) bool
}

type InventoryRequest struct {
	Filter Interface
	Key string
	KeyIsIp bool
	Timeout time.Duration
	Back chan interface{}
}

func NewInvReq(filter Interface, tout time.Duration) InventoryRequest {
	return InventoryRequest{
		Filter: filter,
		Timeout: tout,
		Back: make(chan interface{}, 100),
	}
}

