package parser

import (
	"reflect"
	"sync"
)

var (
	OBJECT_TYPE  = NewClassNode("java.lang.Object", ACC_PUBLIC, nil)
	VOID_TYPE    = NewClassNode("void", 0, nil)
	BOOLEAN_TYPE = NewClassNode("boolean", 0, nil)
	CHAR_TYPE    = NewClassNode("char", 0, nil)
	BYTE_TYPE    = NewClassNode("byte", 0, nil)
	SHORT_TYPE   = NewClassNode("short", 0, nil)
	INT_TYPE     = NewClassNode("int", 0, nil)
	LONG_TYPE    = NewClassNode("long", 0, nil)
	FLOAT_TYPE   = NewClassNode("float", 0, nil)
	DOUBLE_TYPE  = NewClassNode("double", 0, nil)
	STRING_TYPE  = NewClassNode("java.lang.String", ACC_PUBLIC, OBJECT_TYPE)
	LIST_TYPE    = NewClassNode("java.util.List", ACC_PUBLIC, OBJECT_TYPE)
	// New additions
	SCRIPT_TYPE     = NewClassNode("groovy.lang.Script", ACC_PUBLIC, OBJECT_TYPE)
	GSTRING_TYPE    = NewClassNode("groovy.lang.GString", ACC_PUBLIC, OBJECT_TYPE)
	CLOSURE_TYPE    = NewClassNode("groovy.lang.Closure", ACC_PUBLIC, OBJECT_TYPE)
	RANGE_TYPE      = NewClassNode("groovy.lang.Range", ACC_PUBLIC, OBJECT_TYPE)
	PATTERN_TYPE    = NewClassNode("java.util.regex.Pattern", ACC_PUBLIC, OBJECT_TYPE)
	BINDING_TYPE    = NewClassNode("groovy.lang.Binding", ACC_PUBLIC, OBJECT_TYPE)
	BIGINTEGER_TYPE = NewClassNode("java.math.BigInteger", ACC_PUBLIC, OBJECT_TYPE)
	BIGDECIMAL_TYPE = NewClassNode("java.math.BigDecimal", ACC_PUBLIC, OBJECT_TYPE)
	NUMBER_TYPE     = NewClassNode("java.math.Number", ACC_PUBLIC, OBJECT_TYPE)
	MAP_TYPE        = NewClassNode("java.util.Map", ACC_PUBLIC, OBJECT_TYPE)
)

var (
	primitiveClassNames = []string{"boolean", "char", "byte", "short", "int", "long", "float", "double", "void"}
	classes             = []*ClassNode{BOOLEAN_TYPE, CHAR_TYPE, BYTE_TYPE, SHORT_TYPE, INT_TYPE, LONG_TYPE, FLOAT_TYPE, DOUBLE_TYPE, VOID_TYPE}
)

const DYNAMIC_TYPE_METADATA = "_DYNAMIC_TYPE_METADATA_"

func IsPrimitiveVoid(type_ *ClassNode) bool {
	return type_.redirect.Equals(VOID_TYPE)
}

func IsObjectType(type_ *ClassNode) bool {
	return OBJECT_TYPE.Equals(type_)
}

func IsBigIntegerType(type_ *ClassNode) bool {
	return BIGINTEGER_TYPE.Equals(type_)
}

func IsBigDecimalType(type_ *ClassNode) bool {
	return BIGDECIMAL_TYPE.Equals(type_)
}

func dynamicType() *ClassNode {
	node := OBJECT_TYPE.GetPlainNodeReference()
	node.PutNodeMetaData(DYNAMIC_TYPE_METADATA, true)
	return node
}

func MakeFromString(name string) *ClassNode {
	if name == "" {
		return dynamicType()
	}

	for i, primitiveName := range primitiveClassNames {
		if primitiveName == name {
			return classes[i]
		}
	}

	for _, class := range classes {
		if class.GetName() == name {
			return class
		}
	}

	return MakeWithoutCaching(name)
}

func Make(t reflect.Type) *ClassNode {
	switch t.Kind() {
	case reflect.Bool:
		return BOOLEAN_TYPE
	case reflect.Int:
		return INT_TYPE
	case reflect.Int8:
		return BYTE_TYPE
	case reflect.Int16:
		return SHORT_TYPE
	case reflect.Int32:
		return INT_TYPE
	case reflect.Int64:
		return LONG_TYPE
	case reflect.Uint8:
		return BYTE_TYPE
	case reflect.Uint16:
		return SHORT_TYPE
	case reflect.Uint32:
		return INT_TYPE
	case reflect.Uint64:
		return LONG_TYPE
	case reflect.Float32:
		return FLOAT_TYPE
	case reflect.Float64:
		return DOUBLE_TYPE
	case reflect.String:
		return STRING_TYPE
	case reflect.Slice, reflect.Array:
		return LIST_TYPE
	default:
		return NewClassNode(t.String(), ACC_PUBLIC, OBJECT_TYPE)
	}
}

func IsPrimitiveType(cn *ClassNode) bool {
	return cn == BOOLEAN_TYPE || cn == CHAR_TYPE || cn == BYTE_TYPE ||
		cn == SHORT_TYPE || cn == INT_TYPE || cn == LONG_TYPE ||
		cn == FLOAT_TYPE || cn == DOUBLE_TYPE || cn == VOID_TYPE
}

func IsNumberType(cn *ClassNode) bool {
	return cn == BYTE_TYPE || cn == SHORT_TYPE || cn == INT_TYPE ||
		cn == LONG_TYPE || cn == FLOAT_TYPE || cn == DOUBLE_TYPE
}

func IsStringType(cn *ClassNode) bool {
	return cn == STRING_TYPE
}

func IsGStringType(cn *ClassNode) bool {
	return cn == GSTRING_TYPE
}

type ClassHelperCache struct {
	classCache sync.Map // Use sync.Map for concurrent access
}

var globalCache = &ClassHelperCache{}

// MakeCached creates or retrieves a cached ClassNode for the given reflect.Type
func MakeCached(t reflect.Type) *ClassNode {
	// Check if the ClassNode is already in the cache
	if cachedValue, ok := globalCache.classCache.Load(t); ok {
		if classNode, ok := cachedValue.(*ClassNode); ok {
			return classNode
		}
	}

	// If not in cache or invalid, create a new ClassNode
	classNode := NewClassNode(t.Name(), 0, nil) // Adjust parameters as needed

	// Store the new ClassNode in the cache
	globalCache.classCache.Store(t, classNode)

	return classNode
}

func MakeWithoutCaching(name string) *ClassNode {
	cn := NewClassNode(name, ACC_PUBLIC, OBJECT_TYPE)
	cn.SetPrimaryNode(false)
	return cn
}

// New function
func IsDynamicTyped(cn *ClassNode) bool {
	if cn == nil {
		return false
	}
	metadata := cn.GetNodeMetaData(DYNAMIC_TYPE_METADATA)
	return metadata != nil && metadata == true
}

// New additions
var PRIMITIVE_TYPE_TO_WRAPPER_TYPE_MAP = map[*ClassNode]*ClassNode{
	BOOLEAN_TYPE: NewClassNode("java.lang.Boolean", ACC_PUBLIC, OBJECT_TYPE),
	CHAR_TYPE:    NewClassNode("java.lang.Character", ACC_PUBLIC, OBJECT_TYPE),
	BYTE_TYPE:    NewClassNode("java.lang.Byte", ACC_PUBLIC, OBJECT_TYPE),
	SHORT_TYPE:   NewClassNode("java.lang.Short", ACC_PUBLIC, OBJECT_TYPE),
	INT_TYPE:     NewClassNode("java.lang.Integer", ACC_PUBLIC, OBJECT_TYPE),
	LONG_TYPE:    NewClassNode("java.lang.Long", ACC_PUBLIC, OBJECT_TYPE),
	FLOAT_TYPE:   NewClassNode("java.lang.Float", ACC_PUBLIC, OBJECT_TYPE),
	DOUBLE_TYPE:  NewClassNode("java.lang.Double", ACC_PUBLIC, OBJECT_TYPE),
}

func GetWrapper(cn *ClassNode) *ClassNode {
	cn = cn.Redirect()
	if !IsPrimitiveType(cn) {
		return cn
	}

	result, ok := PRIMITIVE_TYPE_TO_WRAPPER_TYPE_MAP[cn]
	if !ok {
		result, ok = PRIMITIVE_TYPE_TO_WRAPPER_TYPE_MAP[cn.Redirect()]
	}

	if result != nil {
		return result
	}

	return cn
}
