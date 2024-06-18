package data

type Directive string

const (
	Method    Directive = "Method"
	Variables Directive = "Variables"
	Test      Directive = "Test"
	Each      Directive = "Each"
	Var       Directive = "Var"
	Url       Directive = "Url"
	RootUrl   Directive = "RootUrl"
	Headers   Directive = "Headers"
	Print     Directive = "Print"
)
