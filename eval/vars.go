package eval

import (
	"encoding/json"
	"fmt"
	"strconv"

	sitter "github.com/smacker/go-tree-sitter"
)

type Value interface {
	isValue()
	json.Marshaler
	json.Unmarshaler
}

func (IntVal) isValue()    {}
func (StringVal) isValue() {}
func (NodeVal) isValue()   {}
func (DictVal) isValue()   {}

type NodeVal struct {
	Src []byte
	N   *sitter.Node
}

func (n *NodeVal) Content() string { return n.N.Content(n.Src) }
func (n *NodeVal) Bytes() []byte   { return n.Src[n.N.StartByte():n.N.EndByte()] }

func (n *NodeVal) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.N.Content(n.Src))
}

func (n *NodeVal) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &n.N)
}

type StringVal struct {
	S string
}

func (s *StringVal) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.S)
}

func (s *StringVal) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.S)
}

type IntVal struct {
	I int
}

func (i *IntVal) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.I)
}

func (i *IntVal) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.I)
}

type DictVal struct {
	D map[Value]Value
}

func (d *DictVal) MarshalJSON() ([]byte, error) {
	var m = make(map[string]Value)
	for k, v := range d.D {
		switch k := k.(type) {
		case *StringVal:
			m[k.S] = v
		case *IntVal:
			m[strconv.Itoa(k.I)] = v
		case *NodeVal:
			m[k.Content()] = v
		default:
			return nil, fmt.Errorf("unsupported key type %T", k)
		}
		m[k.(*StringVal).S] = v
	}
	return json.Marshal(m)
}

func (d *DictVal) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.D)
}
