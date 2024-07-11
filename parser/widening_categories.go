package parser

import (
	"reflect"
	"sort"
	"strings"
)

// WideningCategories provides helper methods to determine the type resulting from
// a widening operation, such as in an addition expression.
type WideningCategories struct{}

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
func (WideningCategories) IsInt(t *ClassNode) bool {
	return t.Equals(INT_TYPE)
}

// IsFloat checks if type is a float
func (WideningCategories) IsFloat(t *ClassNode) bool {
	return t.Equals(FLOAT_TYPE)
}

// IsDouble checks if type is a double
func (WideningCategories) IsDouble(t *ClassNode) bool {
	return t.Equals(DOUBLE_TYPE)
}

// IsIntCategory checks if type is an int, byte, char or short
func (WideningCategories) IsIntCategory(t *ClassNode) bool {
	return t.Equals(INT_TYPE) || t.Equals(BYTE_TYPE) || t.Equals(CHAR_TYPE) || t.Equals(SHORT_TYPE)
}

// IsLongCategory checks if type is a long, int, byte, char or short
func (WideningCategories) IsLongCategory(t *ClassNode) bool {
	return t.Equals(LONG_TYPE) || WideningCategories{}.IsIntCategory(t)
}

// IsBigIntCategory checks if type is a BigInteger, long, int, byte, char or short
func (WideningCategories) IsBigIntCategory(t *ClassNode) bool {
	return IsBigIntegerType(t) || WideningCategories{}.IsLongCategory(t)
}

// IsBigDecCategory checks if type is a BigDecimal, BigInteger, long, int, byte, char or short
func (WideningCategories) IsBigDecCategory(t *ClassNode) bool {
	return IsBigDecimalType(t) || WideningCategories{}.IsBigIntCategory(t)
}

// IsDoubleCategory checks if type is a float, double or BigDecimal (category)
func (WideningCategories) IsDoubleCategory(t *ClassNode) bool {
	return t.Equals(FLOAT_TYPE) || t.Equals(DOUBLE_TYPE) || WideningCategories{}.IsBigDecCategory(t)
}

// IsFloatingCategory checks if type is a float or double
func (WideningCategories) IsFloatingCategory(t *ClassNode) bool {
	return t.Equals(FLOAT_TYPE) || t.Equals(DOUBLE_TYPE)
}

// IsNumberCategory checks if type is a BigDecimal (category) or Number
func (WideningCategories) IsNumberCategory(t *ClassNode) bool {
	return WideningCategories{}.IsBigDecCategory(t) || t.IsDerivedFrom(NUMBER_TYPE)
}

// LowestUpperBound finds the lowest upper bound of a list of types
func (WideningCategories) LowestUpperBound(nodes []*ClassNode) *ClassNode {
	n := len(nodes)
	if n == 1 {
		return nodes[0]
	}
	if n == 2 {
		return WideningCategories{}.lowestUpperBound(nodes[0], nodes[1])
	}
	return WideningCategories{}.lowestUpperBound(nodes[0], WideningCategories{}.LowestUpperBound(nodes[1:]))
}

// lowestUpperBound finds the lowest upper bound of two types
func (WideningCategories) lowestUpperBound(a, b *ClassNode) *ClassNode {
	// This is a simplified version. The actual implementation would be more complex.
	if a.Equals(b) {
		return a
	}
	if IsObjectType(a) || IsObjectType(b) {
		return OBJECT_TYPE
	}
	if a.IsArray() && b.IsArray() {
		return WideningCategories{}.lowestUpperBound(a.GetComponentType(), b.GetComponentType()).MakeArray()
	}
	// More complex logic would be needed here to handle interfaces, inheritance, etc.
	return OBJECT_TYPE
}

// LowestUpperBoundClassNode represents a lowest upper bound that can't be represented by an existing type
type LowestUpperBoundClassNode struct {
	*ClassNode
	upper      *ClassNode
	interfaces []*ClassNode
}

// NewLowestUpperBoundClassNode creates a new LowestUpperBoundClassNode
func NewLowestUpperBoundClassNode(name string, upper *ClassNode, interfaces ...*ClassNode) *LowestUpperBoundClassNode {
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
	var ubs []*ClassNode
	if !IsObjectType(lub.upper) {
		ubs = append(ubs, lub.upper)
	}
	ubs = append(ubs, lub.interfaces...)

	gt := NewGenericsType(MakeWithoutCaching("?"), ubs, nil)
	gt.SetWildcard(true)
	return gt
}

// GetPlainNodeReference returns a plain node reference of the lowest upper bound
func (lub *LowestUpperBoundClassNode) GetPlainNodeReference() *LowestUpperBoundClassNode {
	faces := make([]*ClassNode, len(lub.interfaces))
	for i, iface := range lub.interfaces {
		faces[i] = iface.GetPlainNodeReference()
	}
	return NewLowestUpperBoundClassNode(lub.GetName(), lub.upper.GetPlainNodeReference(), faces...)
}

// GetUpper returns the upper bound of the LowestUpperBoundClassNode
func (lub *LowestUpperBoundClassNode) GetUpper() *ClassNode {
	return lub.upper
}

// GetInterfaces returns the interfaces of the LowestUpperBoundClassNode
func (lub *LowestUpperBoundClassNode) GetInterfaces() []*ClassNode {
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
