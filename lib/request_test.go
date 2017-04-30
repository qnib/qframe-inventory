package qframe_inventory

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)


func TestNewInvReq(t *testing.T) {
	str1 := String{Value: "1"}
	req := NewInvReq(str1, time.Second)
	assert.IsType(t, InventoryRequest{}, req)
}
