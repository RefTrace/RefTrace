// Code generated from GroovyParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // GroovyParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by GroovyParser.
type GroovyParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by GroovyParser#compilationUnit.
	VisitCompilationUnit(ctx *CompilationUnitContext) interface{}

	// Visit a parse tree produced by GroovyParser#scriptStatements.
	VisitScriptStatements(ctx *ScriptStatementsContext) interface{}

	// Visit a parse tree produced by GroovyParser#scriptStatement.
	VisitScriptStatement(ctx *ScriptStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#packageDeclaration.
	VisitPackageDeclaration(ctx *PackageDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#importDeclaration.
	VisitImportDeclaration(ctx *ImportDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeDeclaration.
	VisitTypeDeclaration(ctx *TypeDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#modifier.
	VisitModifier(ctx *ModifierContext) interface{}

	// Visit a parse tree produced by GroovyParser#modifiersOpt.
	VisitModifiersOpt(ctx *ModifiersOptContext) interface{}

	// Visit a parse tree produced by GroovyParser#modifiers.
	VisitModifiers(ctx *ModifiersContext) interface{}

	// Visit a parse tree produced by GroovyParser#classOrInterfaceModifiersOpt.
	VisitClassOrInterfaceModifiersOpt(ctx *ClassOrInterfaceModifiersOptContext) interface{}

	// Visit a parse tree produced by GroovyParser#classOrInterfaceModifiers.
	VisitClassOrInterfaceModifiers(ctx *ClassOrInterfaceModifiersContext) interface{}

	// Visit a parse tree produced by GroovyParser#classOrInterfaceModifier.
	VisitClassOrInterfaceModifier(ctx *ClassOrInterfaceModifierContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableModifier.
	VisitVariableModifier(ctx *VariableModifierContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableModifiersOpt.
	VisitVariableModifiersOpt(ctx *VariableModifiersOptContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableModifiers.
	VisitVariableModifiers(ctx *VariableModifiersContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeParameters.
	VisitTypeParameters(ctx *TypeParametersContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeParameter.
	VisitTypeParameter(ctx *TypeParameterContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeBound.
	VisitTypeBound(ctx *TypeBoundContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeList.
	VisitTypeList(ctx *TypeListContext) interface{}

	// Visit a parse tree produced by GroovyParser#classDeclaration.
	VisitClassDeclaration(ctx *ClassDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#classBody.
	VisitClassBody(ctx *ClassBodyContext) interface{}

	// Visit a parse tree produced by GroovyParser#enumConstants.
	VisitEnumConstants(ctx *EnumConstantsContext) interface{}

	// Visit a parse tree produced by GroovyParser#enumConstant.
	VisitEnumConstant(ctx *EnumConstantContext) interface{}

	// Visit a parse tree produced by GroovyParser#classBodyDeclaration.
	VisitClassBodyDeclaration(ctx *ClassBodyDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#memberDeclaration.
	VisitMemberDeclaration(ctx *MemberDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#methodDeclaration.
	VisitMethodDeclaration(ctx *MethodDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#compactConstructorDeclaration.
	VisitCompactConstructorDeclaration(ctx *CompactConstructorDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#methodName.
	VisitMethodName(ctx *MethodNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#returnType.
	VisitReturnType(ctx *ReturnTypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#fieldDeclaration.
	VisitFieldDeclaration(ctx *FieldDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableDeclarators.
	VisitVariableDeclarators(ctx *VariableDeclaratorsContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableDeclarator.
	VisitVariableDeclarator(ctx *VariableDeclaratorContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableDeclaratorId.
	VisitVariableDeclaratorId(ctx *VariableDeclaratorIdContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableInitializer.
	VisitVariableInitializer(ctx *VariableInitializerContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableInitializers.
	VisitVariableInitializers(ctx *VariableInitializersContext) interface{}

	// Visit a parse tree produced by GroovyParser#emptyDims.
	VisitEmptyDims(ctx *EmptyDimsContext) interface{}

	// Visit a parse tree produced by GroovyParser#emptyDimsOpt.
	VisitEmptyDimsOpt(ctx *EmptyDimsOptContext) interface{}

	// Visit a parse tree produced by GroovyParser#standardType.
	VisitStandardType(ctx *StandardTypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#type.
	VisitType(ctx *TypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#generalClassOrInterfaceType.
	VisitGeneralClassOrInterfaceType(ctx *GeneralClassOrInterfaceTypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#standardClassOrInterfaceType.
	VisitStandardClassOrInterfaceType(ctx *StandardClassOrInterfaceTypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#primitiveType.
	VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeArguments.
	VisitTypeArguments(ctx *TypeArgumentsContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeArgument.
	VisitTypeArgument(ctx *TypeArgumentContext) interface{}

	// Visit a parse tree produced by GroovyParser#annotatedQualifiedClassName.
	VisitAnnotatedQualifiedClassName(ctx *AnnotatedQualifiedClassNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#qualifiedClassNameList.
	VisitQualifiedClassNameList(ctx *QualifiedClassNameListContext) interface{}

	// Visit a parse tree produced by GroovyParser#formalParameters.
	VisitFormalParameters(ctx *FormalParametersContext) interface{}

	// Visit a parse tree produced by GroovyParser#formalParameterList.
	VisitFormalParameterList(ctx *FormalParameterListContext) interface{}

	// Visit a parse tree produced by GroovyParser#thisFormalParameter.
	VisitThisFormalParameter(ctx *ThisFormalParameterContext) interface{}

	// Visit a parse tree produced by GroovyParser#formalParameter.
	VisitFormalParameter(ctx *FormalParameterContext) interface{}

	// Visit a parse tree produced by GroovyParser#methodBody.
	VisitMethodBody(ctx *MethodBodyContext) interface{}

	// Visit a parse tree produced by GroovyParser#qualifiedName.
	VisitQualifiedName(ctx *QualifiedNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#qualifiedNameElement.
	VisitQualifiedNameElement(ctx *QualifiedNameElementContext) interface{}

	// Visit a parse tree produced by GroovyParser#qualifiedNameElements.
	VisitQualifiedNameElements(ctx *QualifiedNameElementsContext) interface{}

	// Visit a parse tree produced by GroovyParser#qualifiedClassName.
	VisitQualifiedClassName(ctx *QualifiedClassNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#qualifiedStandardClassName.
	VisitQualifiedStandardClassName(ctx *QualifiedStandardClassNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#integerLiteralAlt.
	VisitIntegerLiteralAlt(ctx *IntegerLiteralAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#floatingPointLiteralAlt.
	VisitFloatingPointLiteralAlt(ctx *FloatingPointLiteralAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#stringLiteralAlt.
	VisitStringLiteralAlt(ctx *StringLiteralAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#booleanLiteralAlt.
	VisitBooleanLiteralAlt(ctx *BooleanLiteralAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#nullLiteralAlt.
	VisitNullLiteralAlt(ctx *NullLiteralAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#gstring.
	VisitGstring(ctx *GstringContext) interface{}

	// Visit a parse tree produced by GroovyParser#gstringValue.
	VisitGstringValue(ctx *GstringValueContext) interface{}

	// Visit a parse tree produced by GroovyParser#gstringPath.
	VisitGstringPath(ctx *GstringPathContext) interface{}

	// Visit a parse tree produced by GroovyParser#lambdaExpression.
	VisitLambdaExpression(ctx *LambdaExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#standardLambdaExpression.
	VisitStandardLambdaExpression(ctx *StandardLambdaExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#lambdaParameters.
	VisitLambdaParameters(ctx *LambdaParametersContext) interface{}

	// Visit a parse tree produced by GroovyParser#standardLambdaParameters.
	VisitStandardLambdaParameters(ctx *StandardLambdaParametersContext) interface{}

	// Visit a parse tree produced by GroovyParser#lambdaBody.
	VisitLambdaBody(ctx *LambdaBodyContext) interface{}

	// Visit a parse tree produced by GroovyParser#closure.
	VisitClosure(ctx *ClosureContext) interface{}

	// Visit a parse tree produced by GroovyParser#closureOrLambdaExpression.
	VisitClosureOrLambdaExpression(ctx *ClosureOrLambdaExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#blockStatementsOpt.
	VisitBlockStatementsOpt(ctx *BlockStatementsOptContext) interface{}

	// Visit a parse tree produced by GroovyParser#blockStatements.
	VisitBlockStatements(ctx *BlockStatementsContext) interface{}

	// Visit a parse tree produced by GroovyParser#annotationsOpt.
	VisitAnnotationsOpt(ctx *AnnotationsOptContext) interface{}

	// Visit a parse tree produced by GroovyParser#annotation.
	VisitAnnotation(ctx *AnnotationContext) interface{}

	// Visit a parse tree produced by GroovyParser#elementValues.
	VisitElementValues(ctx *ElementValuesContext) interface{}

	// Visit a parse tree produced by GroovyParser#annotationName.
	VisitAnnotationName(ctx *AnnotationNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#elementValuePairs.
	VisitElementValuePairs(ctx *ElementValuePairsContext) interface{}

	// Visit a parse tree produced by GroovyParser#elementValuePair.
	VisitElementValuePair(ctx *ElementValuePairContext) interface{}

	// Visit a parse tree produced by GroovyParser#elementValuePairName.
	VisitElementValuePairName(ctx *ElementValuePairNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#elementValue.
	VisitElementValue(ctx *ElementValueContext) interface{}

	// Visit a parse tree produced by GroovyParser#elementValueArrayInitializer.
	VisitElementValueArrayInitializer(ctx *ElementValueArrayInitializerContext) interface{}

	// Visit a parse tree produced by GroovyParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by GroovyParser#blockStatement.
	VisitBlockStatement(ctx *BlockStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#localVariableDeclaration.
	VisitLocalVariableDeclaration(ctx *LocalVariableDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableDeclaration.
	VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeNamePairs.
	VisitTypeNamePairs(ctx *TypeNamePairsContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeNamePair.
	VisitTypeNamePair(ctx *TypeNamePairContext) interface{}

	// Visit a parse tree produced by GroovyParser#variableNames.
	VisitVariableNames(ctx *VariableNamesContext) interface{}

	// Visit a parse tree produced by GroovyParser#conditionalStatement.
	VisitConditionalStatement(ctx *ConditionalStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#ifElseStatement.
	VisitIfElseStatement(ctx *IfElseStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#switchStatement.
	VisitSwitchStatement(ctx *SwitchStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#forStmtAlt.
	VisitForStmtAlt(ctx *ForStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#whileStmtAlt.
	VisitWhileStmtAlt(ctx *WhileStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#doWhileStmtAlt.
	VisitDoWhileStmtAlt(ctx *DoWhileStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#continueStatement.
	VisitContinueStatement(ctx *ContinueStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#breakStatement.
	VisitBreakStatement(ctx *BreakStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#yieldStatement.
	VisitYieldStatement(ctx *YieldStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#tryCatchStatement.
	VisitTryCatchStatement(ctx *TryCatchStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#assertStatement.
	VisitAssertStatement(ctx *AssertStatementContext) interface{}

	// Visit a parse tree produced by GroovyParser#blockStmtAlt.
	VisitBlockStmtAlt(ctx *BlockStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#conditionalStmtAlt.
	VisitConditionalStmtAlt(ctx *ConditionalStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#loopStmtAlt.
	VisitLoopStmtAlt(ctx *LoopStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#tryCatchStmtAlt.
	VisitTryCatchStmtAlt(ctx *TryCatchStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#synchronizedStmtAlt.
	VisitSynchronizedStmtAlt(ctx *SynchronizedStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#returnStmtAlt.
	VisitReturnStmtAlt(ctx *ReturnStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#throwStmtAlt.
	VisitThrowStmtAlt(ctx *ThrowStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#breakStmtAlt.
	VisitBreakStmtAlt(ctx *BreakStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#continueStmtAlt.
	VisitContinueStmtAlt(ctx *ContinueStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#yieldStmtAlt.
	VisitYieldStmtAlt(ctx *YieldStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#labeledStmtAlt.
	VisitLabeledStmtAlt(ctx *LabeledStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#assertStmtAlt.
	VisitAssertStmtAlt(ctx *AssertStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#localVariableDeclarationStmtAlt.
	VisitLocalVariableDeclarationStmtAlt(ctx *LocalVariableDeclarationStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#expressionStmtAlt.
	VisitExpressionStmtAlt(ctx *ExpressionStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#emptyStmtAlt.
	VisitEmptyStmtAlt(ctx *EmptyStmtAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#catchClause.
	VisitCatchClause(ctx *CatchClauseContext) interface{}

	// Visit a parse tree produced by GroovyParser#catchType.
	VisitCatchType(ctx *CatchTypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#finallyBlock.
	VisitFinallyBlock(ctx *FinallyBlockContext) interface{}

	// Visit a parse tree produced by GroovyParser#resources.
	VisitResources(ctx *ResourcesContext) interface{}

	// Visit a parse tree produced by GroovyParser#resourceList.
	VisitResourceList(ctx *ResourceListContext) interface{}

	// Visit a parse tree produced by GroovyParser#resource.
	VisitResource(ctx *ResourceContext) interface{}

	// Visit a parse tree produced by GroovyParser#switchBlockStatementGroup.
	VisitSwitchBlockStatementGroup(ctx *SwitchBlockStatementGroupContext) interface{}

	// Visit a parse tree produced by GroovyParser#switchLabel.
	VisitSwitchLabel(ctx *SwitchLabelContext) interface{}

	// Visit a parse tree produced by GroovyParser#forControl.
	VisitForControl(ctx *ForControlContext) interface{}

	// Visit a parse tree produced by GroovyParser#enhancedForControl.
	VisitEnhancedForControl(ctx *EnhancedForControlContext) interface{}

	// Visit a parse tree produced by GroovyParser#classicalForControl.
	VisitClassicalForControl(ctx *ClassicalForControlContext) interface{}

	// Visit a parse tree produced by GroovyParser#forInit.
	VisitForInit(ctx *ForInitContext) interface{}

	// Visit a parse tree produced by GroovyParser#forUpdate.
	VisitForUpdate(ctx *ForUpdateContext) interface{}

	// Visit a parse tree produced by GroovyParser#castParExpression.
	VisitCastParExpression(ctx *CastParExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#parExpression.
	VisitParExpression(ctx *ParExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#expressionInPar.
	VisitExpressionInPar(ctx *ExpressionInParContext) interface{}

	// Visit a parse tree produced by GroovyParser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by GroovyParser#expressionListElement.
	VisitExpressionListElement(ctx *ExpressionListElementContext) interface{}

	// Visit a parse tree produced by GroovyParser#enhancedStatementExpression.
	VisitEnhancedStatementExpression(ctx *EnhancedStatementExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#commandExprAlt.
	VisitCommandExprAlt(ctx *CommandExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#postfixExpression.
	VisitPostfixExpression(ctx *PostfixExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#switchExpression.
	VisitSwitchExpression(ctx *SwitchExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#switchBlockStatementExpressionGroup.
	VisitSwitchBlockStatementExpressionGroup(ctx *SwitchBlockStatementExpressionGroupContext) interface{}

	// Visit a parse tree produced by GroovyParser#switchExpressionLabel.
	VisitSwitchExpressionLabel(ctx *SwitchExpressionLabelContext) interface{}

	// Visit a parse tree produced by GroovyParser#postfixExprAlt.
	VisitPostfixExprAlt(ctx *PostfixExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#unaryNotExprAlt.
	VisitUnaryNotExprAlt(ctx *UnaryNotExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#shiftExprAlt.
	VisitShiftExprAlt(ctx *ShiftExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#implicationExprAlt.
	VisitImplicationExprAlt(ctx *ImplicationExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#castExprAlt.
	VisitCastExprAlt(ctx *CastExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#switchExprAlt.
	VisitSwitchExprAlt(ctx *SwitchExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#multipleAssignmentExprAlt.
	VisitMultipleAssignmentExprAlt(ctx *MultipleAssignmentExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#exclusiveOrExprAlt.
	VisitExclusiveOrExprAlt(ctx *ExclusiveOrExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#additiveExprAlt.
	VisitAdditiveExprAlt(ctx *AdditiveExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#regexExprAlt.
	VisitRegexExprAlt(ctx *RegexExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#conditionalExprAlt.
	VisitConditionalExprAlt(ctx *ConditionalExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#powerExprAlt.
	VisitPowerExprAlt(ctx *PowerExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#relationalExprAlt.
	VisitRelationalExprAlt(ctx *RelationalExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#logicalAndExprAlt.
	VisitLogicalAndExprAlt(ctx *LogicalAndExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#assignmentExprAlt.
	VisitAssignmentExprAlt(ctx *AssignmentExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#unaryAddExprAlt.
	VisitUnaryAddExprAlt(ctx *UnaryAddExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#multiplicativeExprAlt.
	VisitMultiplicativeExprAlt(ctx *MultiplicativeExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#inclusiveOrExprAlt.
	VisitInclusiveOrExprAlt(ctx *InclusiveOrExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#logicalOrExprAlt.
	VisitLogicalOrExprAlt(ctx *LogicalOrExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#equalityExprAlt.
	VisitEqualityExprAlt(ctx *EqualityExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#andExprAlt.
	VisitAndExprAlt(ctx *AndExprAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#castExprAltOperand.
	VisitCastExprAltOperand(ctx *CastExprAltOperandContext) interface{}

	// Visit a parse tree produced by GroovyParser#postfixExprAltOperand.
	VisitPostfixExprAltOperand(ctx *PostfixExprAltOperandContext) interface{}

	// Visit a parse tree produced by GroovyParser#unaryNotExprAltOperand.
	VisitUnaryNotExprAltOperand(ctx *UnaryNotExprAltOperandContext) interface{}

	// Visit a parse tree produced by GroovyParser#unaryAddExprAltOperand.
	VisitUnaryAddExprAltOperand(ctx *UnaryAddExprAltOperandContext) interface{}

	// Visit a parse tree produced by GroovyParser#commandExpression.
	VisitCommandExpression(ctx *CommandExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#commandArgument.
	VisitCommandArgument(ctx *CommandArgumentContext) interface{}

	// Visit a parse tree produced by GroovyParser#pathExpression.
	VisitPathExpression(ctx *PathExpressionContext) interface{}

	// Visit a parse tree produced by GroovyParser#pathElement.
	VisitPathElement(ctx *PathElementContext) interface{}

	// Visit a parse tree produced by GroovyParser#namePart.
	VisitNamePart(ctx *NamePartContext) interface{}

	// Visit a parse tree produced by GroovyParser#dynamicMemberName.
	VisitDynamicMemberName(ctx *DynamicMemberNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#indexPropertyArgs.
	VisitIndexPropertyArgs(ctx *IndexPropertyArgsContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedPropertyArgs.
	VisitNamedPropertyArgs(ctx *NamedPropertyArgsContext) interface{}

	// Visit a parse tree produced by GroovyParser#identifierPrmrAlt.
	VisitIdentifierPrmrAlt(ctx *IdentifierPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#literalPrmrAlt.
	VisitLiteralPrmrAlt(ctx *LiteralPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#gstringPrmrAlt.
	VisitGstringPrmrAlt(ctx *GstringPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#newPrmrAlt.
	VisitNewPrmrAlt(ctx *NewPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#thisPrmrAlt.
	VisitThisPrmrAlt(ctx *ThisPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#superPrmrAlt.
	VisitSuperPrmrAlt(ctx *SuperPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#parenPrmrAlt.
	VisitParenPrmrAlt(ctx *ParenPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#closureOrLambdaExpressionPrmrAlt.
	VisitClosureOrLambdaExpressionPrmrAlt(ctx *ClosureOrLambdaExpressionPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#listPrmrAlt.
	VisitListPrmrAlt(ctx *ListPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#mapPrmrAlt.
	VisitMapPrmrAlt(ctx *MapPrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#builtInTypePrmrAlt.
	VisitBuiltInTypePrmrAlt(ctx *BuiltInTypePrmrAltContext) interface{}

	// Visit a parse tree produced by GroovyParser#identifierPrmrAltNamedPropertyArgPrimary.
	VisitIdentifierPrmrAltNamedPropertyArgPrimary(ctx *IdentifierPrmrAltNamedPropertyArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#literalPrmrAltNamedPropertyArgPrimary.
	VisitLiteralPrmrAltNamedPropertyArgPrimary(ctx *LiteralPrmrAltNamedPropertyArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#gstringPrmrAltNamedPropertyArgPrimary.
	VisitGstringPrmrAltNamedPropertyArgPrimary(ctx *GstringPrmrAltNamedPropertyArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#parenPrmrAltNamedPropertyArgPrimary.
	VisitParenPrmrAltNamedPropertyArgPrimary(ctx *ParenPrmrAltNamedPropertyArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#listPrmrAltNamedPropertyArgPrimary.
	VisitListPrmrAltNamedPropertyArgPrimary(ctx *ListPrmrAltNamedPropertyArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#mapPrmrAltNamedPropertyArgPrimary.
	VisitMapPrmrAltNamedPropertyArgPrimary(ctx *MapPrmrAltNamedPropertyArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedPropertyArgPrimary.
	VisitNamedPropertyArgPrimary(ctx *NamedPropertyArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#identifierPrmrAltNamedArgPrimary.
	VisitIdentifierPrmrAltNamedArgPrimary(ctx *IdentifierPrmrAltNamedArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#literalPrmrAltNamedArgPrimary.
	VisitLiteralPrmrAltNamedArgPrimary(ctx *LiteralPrmrAltNamedArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#gstringPrmrAltNamedArgPrimary.
	VisitGstringPrmrAltNamedArgPrimary(ctx *GstringPrmrAltNamedArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedArgPrimary.
	VisitNamedArgPrimary(ctx *NamedArgPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#identifierPrmrAltCommandPrimary.
	VisitIdentifierPrmrAltCommandPrimary(ctx *IdentifierPrmrAltCommandPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#literalPrmrAltCommandPrimary.
	VisitLiteralPrmrAltCommandPrimary(ctx *LiteralPrmrAltCommandPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#gstringPrmrAltCommandPrimary.
	VisitGstringPrmrAltCommandPrimary(ctx *GstringPrmrAltCommandPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#commandPrimary.
	VisitCommandPrimary(ctx *CommandPrimaryContext) interface{}

	// Visit a parse tree produced by GroovyParser#list.
	VisitList(ctx *ListContext) interface{}

	// Visit a parse tree produced by GroovyParser#map.
	VisitMap(ctx *MapContext) interface{}

	// Visit a parse tree produced by GroovyParser#mapEntryList.
	VisitMapEntryList(ctx *MapEntryListContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedPropertyArgList.
	VisitNamedPropertyArgList(ctx *NamedPropertyArgListContext) interface{}

	// Visit a parse tree produced by GroovyParser#mapEntry.
	VisitMapEntry(ctx *MapEntryContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedPropertyArg.
	VisitNamedPropertyArg(ctx *NamedPropertyArgContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedArg.
	VisitNamedArg(ctx *NamedArgContext) interface{}

	// Visit a parse tree produced by GroovyParser#mapEntryLabel.
	VisitMapEntryLabel(ctx *MapEntryLabelContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedPropertyArgLabel.
	VisitNamedPropertyArgLabel(ctx *NamedPropertyArgLabelContext) interface{}

	// Visit a parse tree produced by GroovyParser#namedArgLabel.
	VisitNamedArgLabel(ctx *NamedArgLabelContext) interface{}

	// Visit a parse tree produced by GroovyParser#creator.
	VisitCreator(ctx *CreatorContext) interface{}

	// Visit a parse tree produced by GroovyParser#dim.
	VisitDim(ctx *DimContext) interface{}

	// Visit a parse tree produced by GroovyParser#arrayInitializer.
	VisitArrayInitializer(ctx *ArrayInitializerContext) interface{}

	// Visit a parse tree produced by GroovyParser#anonymousInnerClassDeclaration.
	VisitAnonymousInnerClassDeclaration(ctx *AnonymousInnerClassDeclarationContext) interface{}

	// Visit a parse tree produced by GroovyParser#createdName.
	VisitCreatedName(ctx *CreatedNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#nonWildcardTypeArguments.
	VisitNonWildcardTypeArguments(ctx *NonWildcardTypeArgumentsContext) interface{}

	// Visit a parse tree produced by GroovyParser#typeArgumentsOrDiamond.
	VisitTypeArgumentsOrDiamond(ctx *TypeArgumentsOrDiamondContext) interface{}

	// Visit a parse tree produced by GroovyParser#arguments.
	VisitArguments(ctx *ArgumentsContext) interface{}

	// Visit a parse tree produced by GroovyParser#argumentList.
	VisitArgumentList(ctx *ArgumentListContext) interface{}

	// Visit a parse tree produced by GroovyParser#enhancedArgumentListInPar.
	VisitEnhancedArgumentListInPar(ctx *EnhancedArgumentListInParContext) interface{}

	// Visit a parse tree produced by GroovyParser#firstArgumentListElement.
	VisitFirstArgumentListElement(ctx *FirstArgumentListElementContext) interface{}

	// Visit a parse tree produced by GroovyParser#argumentListElement.
	VisitArgumentListElement(ctx *ArgumentListElementContext) interface{}

	// Visit a parse tree produced by GroovyParser#enhancedArgumentListElement.
	VisitEnhancedArgumentListElement(ctx *EnhancedArgumentListElementContext) interface{}

	// Visit a parse tree produced by GroovyParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by GroovyParser#className.
	VisitClassName(ctx *ClassNameContext) interface{}

	// Visit a parse tree produced by GroovyParser#identifier.
	VisitIdentifier(ctx *IdentifierContext) interface{}

	// Visit a parse tree produced by GroovyParser#builtInType.
	VisitBuiltInType(ctx *BuiltInTypeContext) interface{}

	// Visit a parse tree produced by GroovyParser#keywords.
	VisitKeywords(ctx *KeywordsContext) interface{}

	// Visit a parse tree produced by GroovyParser#rparen.
	VisitRparen(ctx *RparenContext) interface{}

	// Visit a parse tree produced by GroovyParser#nls.
	VisitNls(ctx *NlsContext) interface{}

	// Visit a parse tree produced by GroovyParser#sep.
	VisitSep(ctx *SepContext) interface{}
}
