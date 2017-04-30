package qframe_inventory

/***** Example Interface String
The Interface used within the Inventory must implement the Euqal() method
 */

type String struct {
	Value string
}

func (s String) Equal(other interface{}) bool {
	ostr := other.(String)
	return s.Value == ostr.Value
}
