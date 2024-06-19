package data

type Directive string

const (
	Method    Directive = "Method"
	Variables Directive = "Variables"
	Test      Directive = "Test"
	Var       Directive = "Var"
	Url       Directive = "Url"
	Headers   Directive = "Headers"
	Print     Directive = "Print"
)
