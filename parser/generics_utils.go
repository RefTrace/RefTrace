package parser

func MakeClassSafe0(type_ IClassNode, genericTypes ...*GenericsType) IClassNode {
	plainNodeReference := NewClass(type_)
	if len(genericTypes) > 0 {
		plainNodeReference.SetGenericsTypes(genericTypes)
		if type_.IsGenericsPlaceHolder() {
			plainNodeReference.SetGenericsPlaceHolder(true)
		}
	}
	return plainNodeReference
}

func CreateGenericsSpec(typ IClassNode) map[string]IClassNode {
	return CreateGenericsSpecWithOldSpec(typ, nil)
}

func CreateGenericsSpecWithOldSpec(typ IClassNode, oldSpec map[string]IClassNode) map[string]IClassNode {
	// Example:
	// abstract class A<X,Y,Z> { ... }
	// class C<T extends Number> extends A<T,Object,String> { }
	// the type "A<T,Object,String> -> A<X,Y,Z>" will produce [X:Number,Y:Object,Z:String]

	oc := typ.GetNodeMetaData("outer.class").(*ClassNode) // GROOVY-10646: outer class type parameters
	var newSpec map[string]IClassNode
	if oc != nil {
		newSpec = CreateGenericsSpecWithOldSpec(oc, oldSpec)
	} else {
		newSpec = make(map[string]IClassNode)
	}

	gt := typ.GetGenericsTypes()
	rgt := typ.Redirect().GetGenericsTypes()
	if gt != nil && rgt != nil {
		for i := 0; i < len(gt) && i < len(rgt); i++ {
			correctedType := CorrectToGenericsSpec(oldSpec, gt[i])
			newSpec[rgt[i].GetName()] = correctedType
		}
	}
	return newSpec
}

func CorrectToGenericsSpec(genericsSpec map[string]IClassNode, typ *GenericsType) IClassNode {
	var cn IClassNode

	if typ.IsPlaceholder() {
		name := typ.GetName()
		if name[0] != '#' {
			if gt, ok := genericsSpec[name]; ok {
				cn = gt
			}
		}
	} else if typ.IsWildcard() {
		upperBounds := typ.GetUpperBounds()
		if len(upperBounds) > 0 {
			cn = upperBounds[0] // GROOVY-9891
		}
	}

	if cn == nil {
		cn = typ.GetType()
	}

	return cn
}

func CorrectToGenericsSpecRecurse(genericsSpec map[string]IClassNode, type_ IClassNode) IClassNode {
	return CorrectToGenericsSpecRecurseWithExclusions(genericsSpec, type_, []string{})
}

func CorrectToGenericsSpecRecurseArray(genericsSpec map[string]IClassNode, types []IClassNode) []IClassNode {
	if len(types) == 1 {
		return types
	}
	newTypes := make([]IClassNode, len(types))
	modified := false
	for i, t := range types {
		newTypes[i] = CorrectToGenericsSpecRecurseWithExclusions(genericsSpec, t, []string{})
		modified = modified || (types[i] != newTypes[i])
	}
	if !modified {
		return types
	}
	return newTypes
}

func MakeClassSafeWithGenerics(type_ IClassNode, genericTypes ...*GenericsType) IClassNode {
	if type_.IsArray() {
		return MakeClassSafeWithGenerics(type_.GetComponentType(), genericTypes...).MakeArray()
	}

	nTypes := len(genericTypes)
	var gTypes []*GenericsType

	if nTypes == 0 {
		gTypes = EmptyGenericsTypeArray
	} else {
		gTypes = make([]*GenericsType, nTypes)
		copy(gTypes, genericTypes)
	}

	return MakeClassSafe0(type_, gTypes...)
}

func NewClass(type_ IClassNode) IClassNode {
	return type_.GetPlainNodeReference()
}

func CorrectToGenericsSpecRecurseWithExclusions(genericsSpec map[string]IClassNode, type_ IClassNode, exclusions []string) IClassNode {
	if type_.IsArray() {
		return CorrectToGenericsSpecRecurseWithExclusions(genericsSpec, type_.GetComponentType(), exclusions).MakeArray()
	}
	name := type_.GetUnresolvedName()
	if type_.IsGenericsPlaceHolder() && !containsString(exclusions, name) {
		exclusions = append(exclusions, name) // GROOVY-7722
		if t, ok := genericsSpec[name]; ok {
			type_ = t
			if type_ != nil && type_.IsGenericsPlaceHolder() {
				if type_.GetGenericsTypes() == nil {
					placeholder := MakeWithoutCaching(type_.GetUnresolvedName())
					placeholder.SetGenericsPlaceHolder(true)
					return MakeClassSafeWithGenerics(type_, NewGenericsTypeWithBasicType(placeholder))
				} else if name != type_.GetUnresolvedName() {
					return CorrectToGenericsSpecRecurseWithExclusions(genericsSpec, type_, exclusions)
				}
			}
		}
	}
	if type_ == nil {
		type_ = OBJECT_TYPE.GetPlainNodeReference()
	}
	oldgTypes := type_.GetGenericsTypes()
	var newgTypes []*GenericsType
	if oldgTypes != nil {
		newgTypes = make([]*GenericsType, len(oldgTypes))
		for i, oldgType := range oldgTypes {
			if oldgType.IsWildcard() {
				oldUpper := oldgType.GetUpperBounds()
				var upper []IClassNode
				if oldUpper != nil {
					upper = make([]IClassNode, len(oldUpper))
					for j, u := range oldUpper {
						upper[j] = CorrectToGenericsSpecRecurseWithExclusions(genericsSpec, u, exclusions)
					}
				}
				oldLower := oldgType.GetLowerBound()
				var lower IClassNode
				if oldLower != nil {
					lower = CorrectToGenericsSpecRecurseWithExclusions(genericsSpec, oldLower, exclusions)
				}
				fixed := NewGenericsType(oldgType.GetType(), upper, lower)
				fixed.SetWildcard(true)
				newgTypes[i] = fixed
			} else if oldgType.IsPlaceholder() {
				if t, ok := genericsSpec[oldgType.GetName()]; ok {
					newgTypes[i] = NewGenericsTypeWithBasicType(t)
				} else {
					newgTypes[i] = Erasure(oldgType)
				}
			} else {

				// Create a temporary GenericsType from oldgType.GetType()
				tempGT := NewGenericsTypeWithBasicType(oldgType.GetType())
				correctedType := CorrectToGenericsSpec(genericsSpec, tempGT)
				newgTypes[i] = NewGenericsTypeWithBasicType(CorrectToGenericsSpecRecurseWithExclusions(genericsSpec, correctedType, exclusions))
			}
		}
	}
	return MakeClassSafeWithGenerics(type_, newgTypes...)
}

func Erasure(gt *GenericsType) *GenericsType {
	var cn IClassNode
	cn = gt.GetType().Redirect() // discard the placeholder

	if gt.GetType().GetGenericsTypes() != nil {
		gt = gt.GetType().GetGenericsTypes()[0]
	}

	if gt.GetUpperBounds() != nil {
		cn = gt.GetUpperBounds()[0] // TODO: if length > 1 then union type?
	}

	return cn.AsGenericsType()
}

// Helper functions (you might need to implement these)
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func HasUnresolvedGenerics(type_ IClassNode) bool {
	if type_.IsGenericsPlaceHolder() {
		return true
	}
	if type_.IsArray() {
		return HasUnresolvedGenerics(type_.GetComponentType())
	}
	genericsTypes := type_.GetGenericsTypes()
	if genericsTypes != nil {
		for _, genericsType := range genericsTypes {
			if genericsType.IsPlaceholder() {
				return true
			}
			lowerBound := genericsType.GetLowerBound()
			upperBounds := genericsType.GetUpperBounds()
			if lowerBound != nil {
				if HasUnresolvedGenerics(lowerBound) {
					return true
				}
			} else if upperBounds != nil {
				for _, upperBound := range upperBounds {
					if HasUnresolvedGenerics(upperBound) {
						return true
					}
				}
			} else {
				if HasUnresolvedGenerics(genericsType.GetType()) {
					return true
				}
			}
		}
	}
	return false
}

// ... existing code ...

func ExtractPlaceholders(type_ IClassNode) map[*GenericsTypeName]*GenericsType {
	placeholders := make(map[*GenericsTypeName]*GenericsType)
	ExtractPlaceholdersHelper(type_, placeholders)
	return placeholders
}

func ExtractPlaceholdersHelper(type_ IClassNode, placeholders map[*GenericsTypeName]*GenericsType) {
	if type_ == nil {
		return
	}

	if type_.IsArray() {
		ExtractPlaceholdersHelper(type_.GetComponentType(), placeholders)
		return
	}

	if !type_.IsUsingGenerics() || !type_.IsRedirectNode() {
		return
	}

	genericsTypes := type_.GetGenericsTypes()
	if genericsTypes == nil || len(genericsTypes) == 0 {
		return
	}

	// GROOVY-8609, GROOVY-10067, etc.
	if type_.IsGenericsPlaceHolder() {
		gt := genericsTypes[0]
		name := NewGenericsTypeName(gt.GetName())
		if _, exists := placeholders[name]; !exists {
			placeholders[name] = gt
		}
		return
	}

	redirectGenericsTypes := type_.Redirect().GetGenericsTypes()
	if redirectGenericsTypes == nil {
		redirectGenericsTypes = genericsTypes
	} else if len(redirectGenericsTypes) != len(genericsTypes) {
		panic("Expected earlier checking to detect generics parameter arity mismatch")
	}

	typeArguments := make([]*GenericsType, 0, len(genericsTypes))
	for i, rgt := range redirectGenericsTypes {
		if rgt.IsPlaceholder() { // type parameter
			typeArgument := genericsTypes[i]
			name := NewGenericsTypeName(rgt.GetName())
			if _, exists := placeholders[name]; !exists {
				placeholders[name] = typeArgument
				typeArguments = append(typeArguments, typeArgument)
			}
		}
	}

	// examine non-placeholder type args
	for _, gt := range typeArguments {
		if gt.IsWildcard() {
			lowerBound := gt.GetLowerBound()
			if lowerBound != nil {
				ExtractPlaceholdersHelper(lowerBound, placeholders)
			} else {
				upperBounds := gt.GetUpperBounds()
				if upperBounds != nil {
					for _, upperBound := range upperBounds {
						ExtractPlaceholdersHelper(upperBound, placeholders)
					}
				}
			}
		} else if !gt.IsPlaceholder() {
			ExtractPlaceholdersHelper(gt.GetType(), placeholders)
		}
	}
}

// BuildWildcardType generates a wildcard generic type to be used for checks against
// class nodes. See GenericsType.IsCompatibleWith(IClassNode).
func BuildWildcardType(upperBounds ...IClassNode) *GenericsType {
	gt := NewGenericsType(MakeWithoutCaching("?"), upperBounds, nil)
	gt.SetWildcard(true)
	return gt
}

func NonGeneric(type_ IClassNode) IClassNode {
	dims := 0
	temp := type_
	for temp.IsArray() {
		dims++
		temp = temp.GetComponentType()
	}

	isParameterized := false
	// TODO: implement DecompiledClassNode
	isParameterized = temp.IsUsingGenerics()
	/*
		if dcn, ok := temp.(*DecompiledClassNode); ok {
			isParameterized = dcn.IsParameterized()
		} else {
			isParameterized = temp.IsUsingGenerics()
		}
	*/

	if isParameterized {
		result := temp.GetPlainNodeReference()
		result.SetGenericsTypes(nil)
		result.SetUsingGenerics(false)
		for dims > 0 {
			dims--
			result = result.MakeArray()
		}
		return result
	}

	return type_
}
