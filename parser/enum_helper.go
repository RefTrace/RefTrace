package parser

import (
	"fmt"
)

const (
	FS        = ACC_FINAL | ACC_STATIC
	PUBLIC_FS = ACC_PUBLIC | FS
)

type EnumHelper struct{}

func MakeEnumNode(name string, modifiers int, interfaces []*ClassNode, outerClass *ClassNode) *ClassNode {
	modifiers = modifiers | ACC_FINAL | ACC_ENUM
	var enumClass *ClassNode

	if outerClass == nil {
		enumClass = NewClassNodeWithInterfaces(name, modifiers, nil, interfaces, nil)
	} else {
		name = fmt.Sprintf("%s$%s", outerClass.GetName(), name)
		modifiers |= ACC_STATIC
		innerClass := NewInnerClassNodeWithInterfaces(outerClass, name, modifiers, nil, interfaces, nil)
		enumClass = innerClass.ClassNode
	}

	// set super class and generics info
	// "enum X" -> class X extends Enum<X>
	gt := &GenericsType{typ: enumClass} // Changed Type to typ
	superClass := MakeWithoutCaching("java.lang.Enum")
	superClass.genericsTypes = []*GenericsType{gt}
	enumClass.SetSuperClass(superClass)
	superClass.SetRedirect(ENUM_TYPE)

	return enumClass
}

func AddEnumConstant(enumClass *ClassNode, name string, init Expression) *FieldNode {
	modifiers := PUBLIC_FS | ACC_ENUM

	if init != nil {
		if _, ok := init.(*ListExpression); !ok {
			list := NewListExpression()
			list.AddExpression(init)
			init = list
		}
	}

	fn := NewFieldNode(name, modifiers, enumClass.GetPlainNodeReference(), enumClass, init)
	enumClass.AddField(fn)
	return fn
}
