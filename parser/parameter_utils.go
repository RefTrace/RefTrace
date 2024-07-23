package parser

import (
	"reflect"
)

// IsVargs checks if the last parameter in the array is a varargs parameter
func IsVargs(parameters []*Parameter) bool {
	if len(parameters) == 0 {
		return false
	}
	return parameters[len(parameters)-1].GetType().IsArray()
}

// ParametersEqual checks if two parameter arrays are equal
func ParametersEqual(a, b []*Parameter) bool {
	return parametersEqual(a, b, false)
}

// ParametersEqualWithWrapperType checks if two parameter arrays are equal, considering wrapper types
func ParametersEqualWithWrapperType(a, b []*Parameter) bool {
	return parametersEqual(a, b, true)
}

// ParametersCompatible checks if the source parameters are compatible with the target parameters
func ParametersCompatible(source, target []*Parameter) bool {
	return parametersMatch(source, target, func(sourceType, targetType *ClassNode) bool {
		return IsAssignableTo(sourceType, targetType)
	})
}

// parametersEqual is a helper function to check parameter equality
func parametersEqual(a, b []*Parameter, wrapType bool) bool {
	return parametersMatch(a, b, func(aType, bType *ClassNode) bool {
		if wrapType {
			aType = GetWrapper(aType)
			bType = GetWrapper(bType)
		}
		return reflect.DeepEqual(aType, bType)
	})
}

// parametersMatch is a generic helper function to match parameters based on a type checker function
func parametersMatch(a, b []*Parameter, typeChecker func(*ClassNode, *ClassNode) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		aType := a[i].GetType()
		bType := b[i].GetType()

		if !typeChecker(aType, bType) {
			return false
		}
	}

	return true
}

// IsAssignableTo checks if sourceType is assignable to targetType
// This is a placeholder function and needs to be implemented based on your type system
func IsAssignableTo(sourceType, targetType *ClassNode) bool {
	// Implement the logic to check if sourceType is assignable to targetType
	// This might involve checking inheritance, interfaces, etc.
	return false
}
