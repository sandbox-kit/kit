package model

// SymbolKind distinguishes references whose names may coincide in nested scopes.
type SymbolKind uint8

const (
	ParameterSymbol SymbolKind = iota + 1
	ResultSymbol
	ReceiverSymbol
	CancellationSymbol
	ErrorResultSymbol
	FieldSymbol
	LocalSymbol
)

type SymbolReference struct {
	Kind  SymbolKind
	ID    string
	Owner string
}
