package qframe_inventory

import (
	"errors"
	"sync"
	"fmt"
)

const (
	version = "0.1.0"
)

type Inventory struct {
	Version string
	Data   map[string]Interface
	PendingRequests []InventoryRequest
	mux sync.Mutex
}

func NewInventory() Inventory {
	return Inventory{
		Version: version,
		Data: make(map[string]Interface),
		PendingRequests: []InventoryRequest{},
	}
}

func (i *Inventory) SetItem(key string, item Interface) (err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	i.Data[key] = item
	return
}

func (i *Inventory) GetItem(key string) (cntOut Interface, err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if item, ok := i.Data[key];ok {
		return item, err
	}
	return cntOut, errors.New(fmt.Sprintf("No item found with key '%s'", key))
}

func filterItem(in Interface, filter Interface) (out Interface, err error) {
	if in.Equal(filter) {
		return in, err
	}
	return out, errors.New("filter does not match")
}



func (i *Inventory) HandleRequest(req InventoryRequest) (err error) {
	for _, item := range i.Data {
		res, err := filterItem(item, req.Filter)
		if err == nil {
			req.Back <- res
			return err
		}
	}
	return errors.New("Could not match filter")
}


func (i *Inventory) ServeRequest(req InventoryRequest) {
	err := i.HandleRequest(req)
	if err != nil {
		i.mux.Lock()
		i.PendingRequests = append(i.PendingRequests, req)
		i.mux.Unlock()
	}
}


// CheckRequests iterates over all requests and responses if the request can be fulfilled
func (inv *Inventory) CheckRequests() {
	inv.mux.Lock()
	for i, req := range inv.PendingRequests {
		err := inv.HandleRequest(req)
		if err == nil {
			inv.PendingRequests = append(inv.PendingRequests[:i], inv.PendingRequests[i+1:]...)
		}
	}
	inv.mux.Unlock()
}

/*
func (ci *ContainerInventory) GetCntByEvent(ce qtypes.ContainerEvent) (cnt types.ContainerJSON, err error) {
	//ci.mux.Lock()
	//defer ci.mux.Unlock()
	id := ce.Event.Actor.ID
	cnt = ce.Container
	if ce.Event.Type != "container" {
		return
	}
	switch ce.Event.Action {
	case "die", "destroy":
		if _, ok := ci.IDtoIP[id]; ok {
			delete(ci.IDtoIP, id)
		}
		if _, ok := ci.Data[id]; ok {
			delete(ci.Data, id)
		}
		return
	case "start":
		ci.Data[id] = cnt
		if cnt.State.Running {
			for _, v := range cnt.NetworkSettings.Networks {
				ci.IDtoIP[id] = v.IPAddress
			}
		}
	}
	return cnt, err
}
*/
