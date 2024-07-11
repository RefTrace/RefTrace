package parser

// nullX returns a constant expression representing null
func nullX() Expression {
	return NewConstantExpression(nil)
}
