package qframe_inventory

import (
	"errors"
	"sync"
	"fmt"
	"github.com/docker/docker/api/types"

)

const (
	version = "0.2.0"
)

type Inventory struct {
	Version string
	Data   map[string]types.ContainerJSON
	PendingRequests []ContainerRequest
	mux sync.Mutex
}

func NewInventory() Inventory {
	return Inventory{
		Version: version,
		Data: make(map[string]types.ContainerJSON),
		PendingRequests: []ContainerRequest{},
	}
}

func (i *Inventory) SetItem(key string, item types.ContainerJSON) (err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	i.Data[key] = item
	return
}

func (i *Inventory) GetItem(key string) (cntOut types.ContainerJSON, err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if item, ok := i.Data[key];ok {
		return item, err
	}
	return cntOut, errors.New(fmt.Sprintf("No item found with key '%s'", key))
}

func filterItem(in ContainerRequest, other types.ContainerJSON) (out types.ContainerJSON, err error) {
	fmt.Println(" >> in filterItem)")
	if in.Equal(other) {
		fmt.Println("is True!")
		return other, err
	}
	return out, errors.New("filter does not match")
}



func (i *Inventory) HandleRequest(req ContainerRequest) (err error) {
	if len(i.Data) == 0 {
		return errors.New("Inventory is empty so far")
	}
	for _, cnt := range i.Data {
		res, err := filterItem(req, cnt)
		if err == nil {
			req.Back <- res
			return err
		}
	}
	return errors.New("Could not match filter")
}


func (i *Inventory) ServeRequest(req ContainerRequest) {
	err := i.HandleRequest(req)
	if err != nil {
		fmt.Println("  > Append req to list")
		i.mux.Lock()
		i.PendingRequests = append(i.PendingRequests, req)
		i.mux.Unlock()
	} else {
		fmt.Println("HandleRequest sucessful")
	}
}


// CheckRequests iterates over all requests and responses if the request can be fulfilled
func (inv *Inventory) CheckRequests() {
	if len(inv.PendingRequests) == 0 {
		return
	}
	inv.mux.Lock()
	defer inv.mux.Unlock()
	for i, req := range inv.PendingRequests {
		err := inv.HandleRequest(req)
		if err != nil {
			fmt.Printf("  > %s: %v\n", i, err.Error())
		} else {
			fmt.Printf("  > %s: OK\n", i)
		}
		/*if err == nil {
			fmt.Println(" >> HandleRegest was sucessful")
			if len(inv.PendingRequests) == 1 {
				inv.PendingRequests = []InventoryRequest{}
			} else {
				inv.PendingRequests = append(inv.PendingRequests[:i], inv.PendingRequests[i+1:]...)
			}
		}
		*/
	}
}

