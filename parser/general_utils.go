package parser

// nullX returns a constant expression representing null
func nullX() Expression {
	return NewConstantExpression(nil)
}

func ClassX(clazz *ClassNode) *ClassExpression {
	return NewClassExpression(clazz)
}

func IsOrImplements(type_ IClassNode, interfaceType IClassNode) bool {
	return type_.Equals(interfaceType) || type_.ImplementsInterface(interfaceType)
}
