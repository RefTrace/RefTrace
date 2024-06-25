// Code generated from GroovyParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // GroovyParser
import "github.com/antlr4-go/antlr/v4"

type BaseGroovyParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseGroovyParserVisitor) VisitCompilationUnit(ctx *CompilationUnitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitScriptStatements(ctx *ScriptStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitScriptStatement(ctx *ScriptStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPackageDeclaration(ctx *PackageDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitImportDeclaration(ctx *ImportDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeDeclaration(ctx *TypeDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitModifier(ctx *ModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitModifiersOpt(ctx *ModifiersOptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitModifiers(ctx *ModifiersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassOrInterfaceModifiersOpt(ctx *ClassOrInterfaceModifiersOptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassOrInterfaceModifiers(ctx *ClassOrInterfaceModifiersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassOrInterfaceModifier(ctx *ClassOrInterfaceModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableModifier(ctx *VariableModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableModifiersOpt(ctx *VariableModifiersOptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableModifiers(ctx *VariableModifiersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeParameters(ctx *TypeParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeParameter(ctx *TypeParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeBound(ctx *TypeBoundContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeList(ctx *TypeListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassDeclaration(ctx *ClassDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassBody(ctx *ClassBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEnumConstants(ctx *EnumConstantsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEnumConstant(ctx *EnumConstantContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassBodyDeclaration(ctx *ClassBodyDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMemberDeclaration(ctx *MemberDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMethodDeclaration(ctx *MethodDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCompactConstructorDeclaration(ctx *CompactConstructorDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMethodName(ctx *MethodNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitReturnType(ctx *ReturnTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitFieldDeclaration(ctx *FieldDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableDeclarators(ctx *VariableDeclaratorsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableDeclarator(ctx *VariableDeclaratorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableDeclaratorId(ctx *VariableDeclaratorIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableInitializer(ctx *VariableInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableInitializers(ctx *VariableInitializersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEmptyDims(ctx *EmptyDimsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEmptyDimsOpt(ctx *EmptyDimsOptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitStandardType(ctx *StandardTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitType(ctx *TypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGeneralClassOrInterfaceType(ctx *GeneralClassOrInterfaceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitStandardClassOrInterfaceType(ctx *StandardClassOrInterfaceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeArguments(ctx *TypeArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeArgument(ctx *TypeArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAnnotatedQualifiedClassName(ctx *AnnotatedQualifiedClassNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitQualifiedClassNameList(ctx *QualifiedClassNameListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitFormalParameters(ctx *FormalParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitFormalParameterList(ctx *FormalParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitThisFormalParameter(ctx *ThisFormalParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitFormalParameter(ctx *FormalParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMethodBody(ctx *MethodBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitQualifiedName(ctx *QualifiedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitQualifiedNameElement(ctx *QualifiedNameElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitQualifiedNameElements(ctx *QualifiedNameElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitQualifiedClassName(ctx *QualifiedClassNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitQualifiedStandardClassName(ctx *QualifiedStandardClassNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIntegerLiteralAlt(ctx *IntegerLiteralAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitFloatingPointLiteralAlt(ctx *FloatingPointLiteralAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitStringLiteralAlt(ctx *StringLiteralAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBooleanLiteralAlt(ctx *BooleanLiteralAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNullLiteralAlt(ctx *NullLiteralAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGstring(ctx *GstringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGstringValue(ctx *GstringValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGstringPath(ctx *GstringPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLambdaExpression(ctx *LambdaExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitStandardLambdaExpression(ctx *StandardLambdaExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLambdaParameters(ctx *LambdaParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitStandardLambdaParameters(ctx *StandardLambdaParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLambdaBody(ctx *LambdaBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClosure(ctx *ClosureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClosureOrLambdaExpression(ctx *ClosureOrLambdaExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBlockStatementsOpt(ctx *BlockStatementsOptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBlockStatements(ctx *BlockStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAnnotationsOpt(ctx *AnnotationsOptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAnnotation(ctx *AnnotationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitElementValues(ctx *ElementValuesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAnnotationName(ctx *AnnotationNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitElementValuePairs(ctx *ElementValuePairsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitElementValuePair(ctx *ElementValuePairContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitElementValuePairName(ctx *ElementValuePairNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitElementValue(ctx *ElementValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitElementValueArrayInitializer(ctx *ElementValueArrayInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBlockStatement(ctx *BlockStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLocalVariableDeclaration(ctx *LocalVariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeNamePairs(ctx *TypeNamePairsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeNamePair(ctx *TypeNamePairContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitVariableNames(ctx *VariableNamesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitConditionalStatement(ctx *ConditionalStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIfElseStatement(ctx *IfElseStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSwitchStatement(ctx *SwitchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitForStmtAlt(ctx *ForStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitWhileStmtAlt(ctx *WhileStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitDoWhileStmtAlt(ctx *DoWhileStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitContinueStatement(ctx *ContinueStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBreakStatement(ctx *BreakStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitYieldStatement(ctx *YieldStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTryCatchStatement(ctx *TryCatchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAssertStatement(ctx *AssertStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBlockStmtAlt(ctx *BlockStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitConditionalStmtAlt(ctx *ConditionalStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLoopStmtAlt(ctx *LoopStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTryCatchStmtAlt(ctx *TryCatchStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSynchronizedStmtAlt(ctx *SynchronizedStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitReturnStmtAlt(ctx *ReturnStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitThrowStmtAlt(ctx *ThrowStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBreakStmtAlt(ctx *BreakStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitContinueStmtAlt(ctx *ContinueStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitYieldStmtAlt(ctx *YieldStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLabeledStmtAlt(ctx *LabeledStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAssertStmtAlt(ctx *AssertStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLocalVariableDeclarationStmtAlt(ctx *LocalVariableDeclarationStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitExpressionStmtAlt(ctx *ExpressionStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEmptyStmtAlt(ctx *EmptyStmtAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCatchClause(ctx *CatchClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCatchType(ctx *CatchTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitFinallyBlock(ctx *FinallyBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitResources(ctx *ResourcesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitResourceList(ctx *ResourceListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitResource(ctx *ResourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSwitchBlockStatementGroup(ctx *SwitchBlockStatementGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSwitchLabel(ctx *SwitchLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitForControl(ctx *ForControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEnhancedForControl(ctx *EnhancedForControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassicalForControl(ctx *ClassicalForControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitForInit(ctx *ForInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitForUpdate(ctx *ForUpdateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCastParExpression(ctx *CastParExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitParExpression(ctx *ParExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitExpressionInPar(ctx *ExpressionInParContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitExpressionListElement(ctx *ExpressionListElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEnhancedStatementExpression(ctx *EnhancedStatementExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCommandExprAlt(ctx *CommandExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPostfixExpression(ctx *PostfixExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSwitchExpression(ctx *SwitchExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSwitchBlockStatementExpressionGroup(ctx *SwitchBlockStatementExpressionGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSwitchExpressionLabel(ctx *SwitchExpressionLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPostfixExprAlt(ctx *PostfixExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitUnaryNotExprAlt(ctx *UnaryNotExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitShiftExprAlt(ctx *ShiftExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitImplicationExprAlt(ctx *ImplicationExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCastExprAlt(ctx *CastExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSwitchExprAlt(ctx *SwitchExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMultipleAssignmentExprAlt(ctx *MultipleAssignmentExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitExclusiveOrExprAlt(ctx *ExclusiveOrExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAdditiveExprAlt(ctx *AdditiveExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitRegexExprAlt(ctx *RegexExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitConditionalExprAlt(ctx *ConditionalExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPowerExprAlt(ctx *PowerExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitRelationalExprAlt(ctx *RelationalExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLogicalAndExprAlt(ctx *LogicalAndExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAssignmentExprAlt(ctx *AssignmentExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitUnaryAddExprAlt(ctx *UnaryAddExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMultiplicativeExprAlt(ctx *MultiplicativeExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitInclusiveOrExprAlt(ctx *InclusiveOrExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLogicalOrExprAlt(ctx *LogicalOrExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEqualityExprAlt(ctx *EqualityExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAndExprAlt(ctx *AndExprAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCastExprAltOperand(ctx *CastExprAltOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPostfixExprAltOperand(ctx *PostfixExprAltOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitUnaryNotExprAltOperand(ctx *UnaryNotExprAltOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitUnaryAddExprAltOperand(ctx *UnaryAddExprAltOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCommandExpression(ctx *CommandExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCommandArgument(ctx *CommandArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPathExpression(ctx *PathExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitPathElement(ctx *PathElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNamePart(ctx *NamePartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitDynamicMemberName(ctx *DynamicMemberNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIndexPropertyArgs(ctx *IndexPropertyArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNamedPropertyArgs(ctx *NamedPropertyArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIdentifierPrmrAlt(ctx *IdentifierPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLiteralPrmrAlt(ctx *LiteralPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGstringPrmrAlt(ctx *GstringPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNewPrmrAlt(ctx *NewPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitThisPrmrAlt(ctx *ThisPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSuperPrmrAlt(ctx *SuperPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitParenPrmrAlt(ctx *ParenPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClosureOrLambdaExpressionPrmrAlt(ctx *ClosureOrLambdaExpressionPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitListPrmrAlt(ctx *ListPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMapPrmrAlt(ctx *MapPrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBuiltInTypePrmrAlt(ctx *BuiltInTypePrmrAltContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIdentifierPrmrAltNamedPropertyArgPrimary(ctx *IdentifierPrmrAltNamedPropertyArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLiteralPrmrAltNamedPropertyArgPrimary(ctx *LiteralPrmrAltNamedPropertyArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGstringPrmrAltNamedPropertyArgPrimary(ctx *GstringPrmrAltNamedPropertyArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitParenPrmrAltNamedPropertyArgPrimary(ctx *ParenPrmrAltNamedPropertyArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitListPrmrAltNamedPropertyArgPrimary(ctx *ListPrmrAltNamedPropertyArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMapPrmrAltNamedPropertyArgPrimary(ctx *MapPrmrAltNamedPropertyArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIdentifierPrmrAltNamedArgPrimary(ctx *IdentifierPrmrAltNamedArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLiteralPrmrAltNamedArgPrimary(ctx *LiteralPrmrAltNamedArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGstringPrmrAltNamedArgPrimary(ctx *GstringPrmrAltNamedArgPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIdentifierPrmrAltCommandPrimary(ctx *IdentifierPrmrAltCommandPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitLiteralPrmrAltCommandPrimary(ctx *LiteralPrmrAltCommandPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitGstringPrmrAltCommandPrimary(ctx *GstringPrmrAltCommandPrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitList(ctx *ListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMap(ctx *MapContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMapEntryList(ctx *MapEntryListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNamedPropertyArgList(ctx *NamedPropertyArgListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMapEntry(ctx *MapEntryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNamedPropertyArg(ctx *NamedPropertyArgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNamedArg(ctx *NamedArgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitMapEntryLabel(ctx *MapEntryLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNamedPropertyArgLabel(ctx *NamedPropertyArgLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNamedArgLabel(ctx *NamedArgLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCreator(ctx *CreatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitDim(ctx *DimContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitArrayInitializer(ctx *ArrayInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitAnonymousInnerClassDeclaration(ctx *AnonymousInnerClassDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitCreatedName(ctx *CreatedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNonWildcardTypeArguments(ctx *NonWildcardTypeArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitTypeArgumentsOrDiamond(ctx *TypeArgumentsOrDiamondContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitArgumentList(ctx *ArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEnhancedArgumentListInPar(ctx *EnhancedArgumentListInParContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitFirstArgumentListElement(ctx *FirstArgumentListElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitArgumentListElement(ctx *ArgumentListElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitEnhancedArgumentListElement(ctx *EnhancedArgumentListElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitClassName(ctx *ClassNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitIdentifier(ctx *IdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitBuiltInType(ctx *BuiltInTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitKeywords(ctx *KeywordsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitRparen(ctx *RparenContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitNls(ctx *NlsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseGroovyParserVisitor) VisitSep(ctx *SepContext) interface{} {
	return v.VisitChildren(ctx)
}
