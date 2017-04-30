package qframe_inventory

/******************** Inventory Request
 Sends a query for a key or an IP and provides a back-channel, so that the requesting partner can block on the request
 until it arrives - honouring a timeout...
*/

import (
	"time"
	"github.com/docker/docker/api/types"

	"strings"
)

type ContainerRequest struct {
	IssuedAt time.Time
	Name string
	ID string
	IP string
	Back chan types.ContainerJSON
}

func NewContainerRequest() ContainerRequest {
	cr := ContainerRequest{
		IssuedAt: time.Now(),
		Back: make(chan types.ContainerJSON, 5),
	}
	return cr
}


func NewIDContainerRequest(id string) ContainerRequest {
	cr := ContainerRequest{
		ID: id,
		IssuedAt: time.Now(),
		Back: make(chan types.ContainerJSON, 5),
	}
	return cr
}

func NewNameContainerRequest(name string) ContainerRequest {
	cr := ContainerRequest{
		Name: name,
		IssuedAt: time.Now(),
		Back: make(chan types.ContainerJSON, 5),
	}
	return cr
}

func NewIPContainerRequest(ip string) ContainerRequest {
	cr := ContainerRequest{
		IP: ip,
		IssuedAt: time.Now(),
		Back: make(chan types.ContainerJSON, 5),
	}
	return cr
}

func (this ContainerRequest) Equal(other types.ContainerJSON) bool {
	matchIP := false
	if other.NetworkSettings.Networks != nil {
		for _, net := range other.NetworkSettings.Networks {
			if this.IP == net.IPAddress {
				matchIP = true
			}
		}
	}
	return this.ID == other.ID || this.Name == strings.Trim(other.Name, "/") || matchIP
}

