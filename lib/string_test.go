package qframe_inventory

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestString_Equal(t *testing.T) {
	str1 := String{Value: "1"}
	str2 := String{Value: "1"}
	assert.True(t, str1.Equal(str2))
}
