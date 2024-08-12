package parser

import (
	"reflect"
	"sort"
	"strings"
)

var _ IClassNode = (*LowestUpperBoundClassNode)(nil)

// NumberTypesPrecedence maps types to their precedence
var NumberTypesPrecedence = map[reflect.Kind]int{
	reflect.Float64: 0,
	reflect.Float32: 1,
	reflect.Int64:   2,
	reflect.Int32:   3,
	reflect.Int16:   4,
	reflect.Int8:    5,
}

// IsInt checks if type is an int
func IsInt(t *ClassNode) bool {
	return t.Equals(INT_TYPE)
}

// IsFloat checks if type is a float
func IsFloat(t *ClassNode) bool {
	return t.Equals(FLOAT_TYPE)
}

// IsDouble checks if type is a double
func IsDouble(t *ClassNode) bool {
	return t.Equals(DOUBLE_TYPE)
}

// IsIntCategory checks if type is an int, byte, char or short
func IsIntCategory(t *ClassNode) bool {
	return t.Equals(INT_TYPE) || t.Equals(BYTE_TYPE) || t.Equals(CHAR_TYPE) || t.Equals(SHORT_TYPE)
}

// IsLongCategory checks if type is a long, int, byte, char or short
func IsLongCategory(t *ClassNode) bool {
	return t.Equals(LONG_TYPE) || IsIntCategory(t)
}

// IsBigIntCategory checks if type is a BigInteger, long, int, byte, char or short
func IsBigIntCategory(t *ClassNode) bool {
	return IsBigIntegerType(t) || IsLongCategory(t)
}

// IsBigDecCategory checks if type is a BigDecimal, BigInteger, long, int, byte, char or short
func IsBigDecCategory(t *ClassNode) bool {
	return IsBigDecimalType(t) || IsBigIntCategory(t)
}

// IsDoubleCategory checks if type is a float, double or BigDecimal (category)
func IsDoubleCategory(t *ClassNode) bool {
	return t.Equals(FLOAT_TYPE) || t.Equals(DOUBLE_TYPE) || IsBigDecCategory(t)
}

// IsFloatingCategory checks if type is a float or double
func IsFloatingCategory(t IClassNode) bool {
	return t.Equals(FLOAT_TYPE) || t.Equals(DOUBLE_TYPE)
}

// IsNumberCategory checks if type is a BigDecimal (category) or Number
func IsNumberCategory(t *ClassNode) bool {
	return IsBigDecCategory(t) || t.IsDerivedFrom(NUMBER_TYPE)
}

// LowestUpperBound finds the lowest upper bound of a list of types
func LowestUpperBound(nodes []IClassNode) IClassNode {
	n := len(nodes)
	if n == 1 {
		return nodes[0]
	}
	if n == 2 {
		return LowestUpperBoundPair(nodes[0], nodes[1])
	}
	return LowestUpperBoundPair(nodes[0], LowestUpperBound(nodes[1:]))
}

// LowestUpperBoundPair finds the lowest upper bound of two types
func LowestUpperBoundPair(a, b IClassNode) IClassNode {
	// This is a simplified version. The actual implementation would be more complex.
	if a.Equals(b) {
		return a
	}
	if IsObjectType(a) || IsObjectType(b) {
		return OBJECT_TYPE
	}
	if a.IsArray() && b.IsArray() {
		return LowestUpperBoundPair(a.GetComponentType(), b.GetComponentType()).MakeArray()
	}
	// More complex logic would be needed here to handle interfaces, inheritance, etc.
	return OBJECT_TYPE
}

// LowestUpperBoundClassNode represents a lowest upper bound that can't be represented by an existing type
type LowestUpperBoundClassNode struct {
	*ClassNode
	upper      IClassNode
	interfaces []IClassNode
}

// NewLowestUpperBoundClassNode creates a new LowestUpperBoundClassNode
func NewLowestUpperBoundClassNode(name string, upper IClassNode, interfaces ...IClassNode) *LowestUpperBoundClassNode {
	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].GetName() < interfaces[j].GetName()
	})

	lub := &LowestUpperBoundClassNode{
		ClassNode:  NewClassNode(name, ACC_PUBLIC|ACC_FINAL, upper),
		upper:      upper,
		interfaces: interfaces,
	}

	var parts []string
	if !IsObjectType(upper) {
		parts = append(parts, upper.GetName())
	}
	for _, iface := range interfaces {
		parts = append(parts, iface.GetName())
	}
	lub.ClassNode.name = "(" + strings.Join(parts, " & ") + ")"

	return lub
}

// GetLubName returns the name of the lowest upper bound
func (lub *LowestUpperBoundClassNode) GetLubName() string {
	return lub.GetName()
}

// GetText returns the textual representation of the lowest upper bound
func (lub *LowestUpperBoundClassNode) GetText() string {
	return lub.name
}

// AsGenericsType returns the lowest upper bound as a generics type
func (lub *LowestUpperBoundClassNode) AsGenericsType() *GenericsType {
	var ubs []IClassNode
	if !IsObjectType(lub.upper) {
		ubs = append(ubs, lub.upper)
	}
	ubs = append(ubs, lub.interfaces...)

	gt := NewGenericsType(MakeWithoutCaching("?"), ubs, nil)
	gt.SetWildcard(true)
	return gt
}

// GetPlainNodeReference returns a plain node reference of the lowest upper bound
func (lub *LowestUpperBoundClassNode) GetPlainNodeReference() IClassNode {
	faces := make([]IClassNode, len(lub.interfaces))
	for i, iface := range lub.interfaces {
		faces[i] = iface.GetPlainNodeReference()
	}
	return NewLowestUpperBoundClassNode(lub.GetName(), lub.upper.GetPlainNodeReference(), faces...)
}

// GetUpper returns the upper bound of the LowestUpperBoundClassNode
func (lub *LowestUpperBoundClassNode) GetUpper() IClassNode {
	return lub.upper
}

// GetInterfaces returns the interfaces of the LowestUpperBoundClassNode
func (lub *LowestUpperBoundClassNode) GetInterfaces() []IClassNode {
	return lub.interfaces
}

// ImplementsInterfaceOrSubclassOf determines if the source class implements an interface or subclasses the target type
func ImplementsInterfaceOrSubclassOf(source *ClassNode, target interface{}) bool {
	if targetClass, ok := target.(*ClassNode); ok {
		if source.IsDerivedFrom(targetClass) || source.ImplementsInterface(targetClass) {
			return true
		}
	}
	if lub, ok := target.(*LowestUpperBoundClassNode); ok {
		if ImplementsInterfaceOrSubclassOf(source, lub.upper) {
			return true
		}
		for _, iface := range lub.interfaces {
			if source.ImplementsInterface(iface) {
				return true
			}
		}
	}
	return false
}
