package parser

import (
	"reflect"
	"sync"
)

var (
	OBJECT_TYPE                = MakeCached(reflect.TypeOf((*interface{})(nil)).Elem())
	CLOSURE_TYPE               = NewClassNode("groovy.lang.Closure", ACC_PUBLIC, OBJECT_TYPE)
	GSTRING_TYPE               = NewClassNode("groovy.lang.GString", ACC_PUBLIC, OBJECT_TYPE)
	RANGE_TYPE                 = NewClassNode("groovy.lang.Range", ACC_PUBLIC, OBJECT_TYPE)
	PATTERN_TYPE               = NewClassNode("java.util.regex.Pattern", ACC_PUBLIC, OBJECT_TYPE)
	STRING_TYPE                = MakeCached(reflect.TypeOf(""))
	SCRIPT_TYPE                = NewClassNode("groovy.lang.Script", ACC_PUBLIC, OBJECT_TYPE)
	BINDING_TYPE               = NewClassNode("groovy.lang.Binding", ACC_PUBLIC, OBJECT_TYPE)
	THROWABLE_TYPE             = NewClassNode("java.lang.Throwable", ACC_PUBLIC, OBJECT_TYPE)
	BOOLEAN_TYPE               = MakeCached(reflect.TypeOf(false))
	CHAR_TYPE                  = MakeCached(reflect.TypeOf(int32(0))) // Go doesn't have a char type, using rune (int32)
	BYTE_TYPE                  = MakeCached(reflect.TypeOf(int8(0)))
	INT_TYPE                   = MakeCached(reflect.TypeOf(int(0)))
	LONG_TYPE                  = MakeCached(reflect.TypeOf(int64(0)))
	SHORT_TYPE                 = MakeCached(reflect.TypeOf(int16(0)))
	DOUBLE_TYPE                = MakeCached(reflect.TypeOf(float64(0)))
	FLOAT_TYPE                 = MakeCached(reflect.TypeOf(float32(0)))
	VOID_TYPE                  = NewClassNode("void", 0, nil)
	VOID_WRAPPER_TYPE          = NewClassNode("java.lang.Void", ACC_PUBLIC, OBJECT_TYPE)
	METACLASS_TYPE             = NewClassNode("groovy.lang.MetaClass", ACC_PUBLIC, OBJECT_TYPE)
	ITERATOR_TYPE              = NewClassNode("java.util.Iterator", ACC_PUBLIC, OBJECT_TYPE)
	ANNOTATION_TYPE            = NewClassNode("java.lang.annotation.Annotation", ACC_PUBLIC, OBJECT_TYPE)
	ELEMENT_TYPE_TYPE          = NewClassNode("java.lang.annotation.ElementType", ACC_PUBLIC, OBJECT_TYPE)
	AUTOCLOSEABLE_TYPE         = NewClassNode("java.lang.AutoCloseable", ACC_PUBLIC, OBJECT_TYPE)
	SERIALIZABLE_TYPE          = NewClassNode("java.io.Serializable", ACC_PUBLIC, OBJECT_TYPE)
	SERIALIZEDLAMBDA_TYPE      = NewClassNode("java.lang.invoke.SerializedLambda", ACC_PUBLIC, OBJECT_TYPE)
	SEALED_TYPE                = NewClassNode("java.lang.Sealed", ACC_PUBLIC, OBJECT_TYPE)
	OVERRIDE_TYPE              = NewClassNode("java.lang.Override", ACC_PUBLIC, OBJECT_TYPE)
	DEPRECATED_TYPE            = NewClassNode("java.lang.Deprecated", ACC_PUBLIC, OBJECT_TYPE)
	MAP_TYPE                   = MakeWithoutCaching("java.util.Map")
	SET_TYPE                   = MakeWithoutCaching("java.util.Set")
	LIST_TYPE                  = MakeWithoutCaching("java.util.List")
	ENUM_TYPE                  = MakeWithoutCaching("java.lang.Enum")
	CLASS_TYPE                 = MakeWithoutCaching("java.lang.Class")
	TUPLE_TYPE                 = MakeWithoutCaching("groovy.lang.Tuple")
	STREAM_TYPE                = MakeWithoutCaching("java.util.stream.Stream")
	ITERABLE_TYPE              = MakeWithoutCaching("java.lang.Iterable")
	REFERENCE_TYPE             = MakeWithoutCaching("java.lang.ref.Reference")
	COLLECTION_TYPE            = MakeWithoutCaching("java.util.Collection")
	COMPARABLE_TYPE            = MakeWithoutCaching("java.lang.Comparable")
	GROOVY_OBJECT_TYPE         = MakeWithoutCaching("groovy.lang.GroovyObject")
	GENERATED_LAMBDA_TYPE      = MakeWithoutCaching("groovy.lang.GeneratedLambda")
	GENERATED_CLOSURE_TYPE     = MakeWithoutCaching("groovy.lang.GeneratedClosure")
	GROOVY_INTERCEPTABLE_TYPE  = MakeWithoutCaching("groovy.lang.GroovyInterceptable")
	GROOVY_OBJECT_SUPPORT_TYPE = MakeWithoutCaching("groovy.lang.GroovyObjectSupport")
	BIGINTEGER_TYPE            = NewClassNode("java.math.BigInteger", ACC_PUBLIC, OBJECT_TYPE)
	BIGDECIMAL_TYPE            = NewClassNode("java.math.BigDecimal", ACC_PUBLIC, OBJECT_TYPE)
	NUMBER_TYPE                = NewClassNode("java.lang.Number", ACC_PUBLIC, OBJECT_TYPE)

	// Wrapper types for primitives
	BYTE_WRAPPER_TYPE      = NewClassNode("java.lang.Byte", ACC_PUBLIC, OBJECT_TYPE)
	SHORT_WRAPPER_TYPE     = NewClassNode("java.lang.Short", ACC_PUBLIC, OBJECT_TYPE)
	INTEGER_WRAPPER_TYPE   = NewClassNode("java.lang.Integer", ACC_PUBLIC, OBJECT_TYPE)
	LONG_WRAPPER_TYPE      = NewClassNode("java.lang.Long", ACC_PUBLIC, OBJECT_TYPE)
	CHARACTER_WRAPPER_TYPE = NewClassNode("java.lang.Character", ACC_PUBLIC, OBJECT_TYPE)
	FLOAT_WRAPPER_TYPE     = NewClassNode("java.lang.Float", ACC_PUBLIC, OBJECT_TYPE)
	DOUBLE_WRAPPER_TYPE    = NewClassNode("java.lang.Double", ACC_PUBLIC, OBJECT_TYPE)
	BOOLEAN_WRAPPER_TYPE   = NewClassNode("java.lang.Boolean", ACC_PUBLIC, OBJECT_TYPE)
)

var PRIMITIVE_TYPE_TO_DESCRIPTION_MAP = map[IClassNode]string{
	INT_TYPE:     "I",
	VOID_TYPE:    "V",
	BOOLEAN_TYPE: "Z",
	BYTE_TYPE:    "B",
	CHAR_TYPE:    "C",
	SHORT_TYPE:   "S",
	DOUBLE_TYPE:  "D",
	FLOAT_TYPE:   "F",
	LONG_TYPE:    "J",
}

var (
	primitiveClassNames = []string{"boolean", "char", "byte", "short", "int", "long", "float", "double", "void"}
	classes             = []IClassNode{BOOLEAN_TYPE, CHAR_TYPE, BYTE_TYPE, SHORT_TYPE, INT_TYPE, LONG_TYPE, FLOAT_TYPE, DOUBLE_TYPE, VOID_TYPE}
)

const DYNAMIC_TYPE_METADATA = "_DYNAMIC_TYPE_METADATA_"

// TODO: implement traits
func HasDefaultImplementation(method MethodOrConstructorNode) bool {
	// assume all methods have a default implementation in the trait
	return true
}

// findSAM returns the single abstract method of a class node, if it is a SAM type, or nil otherwise.
//
// Parameters:
//   - type_: a type for which to search for a single abstract method
//
// Returns:
//   - the method node if type_ is a SAM type, nil otherwise
func FindSAM(type_ IClassNode) MethodOrConstructorNode {
	if type_ == nil {
		return nil
	}
	if type_.IsInterface() {
		var sam MethodOrConstructorNode
		for _, mn := range type_.GetAbstractMethods() {
			// ignore methods that will have an implementation
			if HasDefaultImplementation(mn) {
				continue
			}
			/*
				name := mn.GetName()
				if OBJECT_METHOD_NAME_SET[name] {
					// Avoid unnecessary checking for `Object` methods as possible as we could
					if OBJECT_TYPE.GetDeclaredMethod(name, mn.GetParameters()) != nil {
						continue
					}
				}

				// we have two methods, so no SAM
				if sam != nil {
					return nil
				}
				sam = mn
			*/
		}
		return sam
	}
	// TODO: implement abstract methods
	/*
		if type_.IsAbstract() {
			var sam MethodOrConstructorNode
			for _, mn := range type_.GetAbstractMethods() {
				if !hasUsableImplementation(type_, mn) {
					if sam != nil {
						return nil
					}
					sam = mn
				}
			}
			return sam
		}
	*/
	return nil
}

func IsPrimitiveVoid(type_ IClassNode) bool {
	return type_.Redirect().Equals(VOID_TYPE)
}

func IsObjectType(type_ IClassNode) bool {
	return OBJECT_TYPE.Equals(type_)
}

func IsBigIntegerType(type_ IClassNode) bool {
	return BIGINTEGER_TYPE.Equals(type_)
}

func IsBigDecimalType(type_ IClassNode) bool {
	return BIGDECIMAL_TYPE.Equals(type_)
}

func IsClassType(type_ IClassNode) bool {
	return CLASS_TYPE.Equals(type_)
}

func DynamicType() IClassNode {
	node := OBJECT_TYPE.GetPlainNodeReference()
	node.PutNodeMetaData(DYNAMIC_TYPE_METADATA, true)
	return node.(IClassNode)
}

func MakeFromString(name string) IClassNode {
	if name == "" {
		return DynamicType()
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

func Make(t reflect.Type) IClassNode {
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
		if t.String() == "tupletype" {
			return TUPLE_TYPE
		}
		return STRING_TYPE
	case reflect.Slice, reflect.Array:
		return LIST_TYPE
	default:
		return NewClassNode(t.String(), ACC_PUBLIC, OBJECT_TYPE)
	}
}

func IsGroovyObjectType(type_ IClassNode) bool {
	return GROOVY_OBJECT_TYPE.Equals(type_)
}

func IsPrimitiveType(cn IClassNode) bool {
	if cn == nil {
		return false
	}
	_, exists := PRIMITIVE_TYPE_TO_DESCRIPTION_MAP[cn.Redirect()]
	return exists
}

func IsNumberType(cn IClassNode) bool {
	return cn == BYTE_TYPE || cn == SHORT_TYPE || cn == INT_TYPE ||
		cn == LONG_TYPE || cn == FLOAT_TYPE || cn == DOUBLE_TYPE
}

func IsStringType(cn IClassNode) bool {
	return cn == STRING_TYPE
}

func IsGStringType(cn IClassNode) bool {
	return cn == GSTRING_TYPE
}

type ClassHelperCache struct {
	classCache sync.Map // Use sync.Map for concurrent access
}

var globalCache = &ClassHelperCache{}

// MakeCached creates or retrieves a cached ClassNode for the given reflect.Type
func MakeCached(t reflect.Type) IClassNode {
	// Check if the ClassNode is already in the cache
	if cachedValue, ok := globalCache.classCache.Load(t); ok {
		if classNode, ok := cachedValue.(IClassNode); ok {
			return classNode
		}
	}

	// If not in cache or invalid, create a new ClassNode
	classNode := NewClassNode(t.Name(), 0, nil) // Adjust parameters as needed

	// Store the new ClassNode in the cache
	globalCache.classCache.Store(t, classNode)

	return classNode
}

func MakeWithoutCaching(name string) IClassNode {
	cn := NewClassNode(name, ACC_PUBLIC, OBJECT_TYPE)
	cn.SetPrimaryNode(false)
	return cn
}

// New function
func IsDynamicTyped(cn IClassNode) bool {
	if cn == nil {
		return false
	}
	metadata := cn.GetNodeMetaData(DYNAMIC_TYPE_METADATA)
	return metadata != nil && metadata == true
}

var WRAPPER_TYPE_TO_PRIMITIVE_TYPE_MAP map[IClassNode]IClassNode

func init() {
	WRAPPER_TYPE_TO_PRIMITIVE_TYPE_MAP = make(map[IClassNode]IClassNode)
	for k, v := range PRIMITIVE_TYPE_TO_WRAPPER_TYPE_MAP {
		WRAPPER_TYPE_TO_PRIMITIVE_TYPE_MAP[v] = k
	}
}

func GetUnwrapper(cn IClassNode) IClassNode {
	cn = cn.Redirect()
	if IsPrimitiveType(cn) {
		return cn
	}

	result, ok := WRAPPER_TYPE_TO_PRIMITIVE_TYPE_MAP[cn]
	if ok {
		return result
	}

	return cn
}

// New additions
var PRIMITIVE_TYPE_TO_WRAPPER_TYPE_MAP = map[IClassNode]IClassNode{
	BOOLEAN_TYPE: NewClassNode("java.lang.Boolean", ACC_PUBLIC, OBJECT_TYPE),
	CHAR_TYPE:    NewClassNode("java.lang.Character", ACC_PUBLIC, OBJECT_TYPE),
	BYTE_TYPE:    NewClassNode("java.lang.Byte", ACC_PUBLIC, OBJECT_TYPE),
	SHORT_TYPE:   NewClassNode("java.lang.Short", ACC_PUBLIC, OBJECT_TYPE),
	INT_TYPE:     NewClassNode("java.lang.Integer", ACC_PUBLIC, OBJECT_TYPE),
	LONG_TYPE:    NewClassNode("java.lang.Long", ACC_PUBLIC, OBJECT_TYPE),
	FLOAT_TYPE:   NewClassNode("java.lang.Float", ACC_PUBLIC, OBJECT_TYPE),
	DOUBLE_TYPE:  NewClassNode("java.lang.Double", ACC_PUBLIC, OBJECT_TYPE),
	VOID_TYPE:    VOID_WRAPPER_TYPE,
}

func GetWrapper(cn IClassNode) IClassNode {
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

func IsSAMType(cn IClassNode) bool {
	// TODO: Implement this
	return false
}
