package qframe_inventory


import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

var (
	str1 = String{Value: "1"}
	str2 = String{Value: "2"}
	str3 = String{Value: "3"}
)

func TestInventory_SetItem(t *testing.T) {
	i := NewInventory()
	assert.IsType(t, Inventory{}, i)
	i.Data["ID1"] = str1
	assert.Len(t, i.Data, 1)
	i.SetItem("ID2", str2)
	assert.Len(t, i.Data, 2)
}

func TestInventory_GetItem(t *testing.T) {
	i := NewInventory()
	i.SetItem("ID1", str1)
	i.SetItem("ID2", str2)
	got, err := i.GetItem("ID1")
	assert.NoError(t, err)
	assert.Equal(t, str1, got)
}

func TestInventory_filterItem(t *testing.T) {
	got, err := filterItem(str1, str1)
	assert.NoError(t, err)
	assert.Equal(t, str1, got)
	got, err = filterItem(str1, str2)
	assert.Error(t, err, err.Error())
}

func TestInventory_HandleRequest(t *testing.T) {
	i := NewInventory()
	i.SetItem("ID1", str1)
	req := NewInvReq(str1, time.Second)
	err := i.HandleRequest(req)
	assert.NoError(t, err)
	m := <-req.Back
	assert.Equal(t, str1, m)
	req = NewInvReq(str2, 5*time.Second)
	err = i.HandleRequest(req)
	assert.Error(t, err)
}

func TestInventory_ServeRequest(t *testing.T) {
	i := NewInventory()
	i.SetItem("ID1", str1)
	req := NewInvReq(str1, 5*time.Second)
	i.ServeRequest(req)
	assert.Equal(t, len(i.PendingRequests), 0)
	req = NewInvReq(str2, 5*time.Second)
	i.ServeRequest(req)
	assert.Equal(t, len(i.PendingRequests), 1)
}

func TestInventory_CheckRequest(t *testing.T) {
	i := NewInventory()
	i.SetItem("ID1", str1)
	req := NewInvReq(str2, 5*time.Second)
	i.ServeRequest(req)
	i.SetItem("ID2", str2)
	i.CheckRequests()
	m := <-req.Back
	assert.Equal(t, str2, m)
}
