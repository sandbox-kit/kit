package declarationgen

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"go/token"
)

// Symbol resolves an explicit reference in the current declaration. Local names
// remain independent of parameter names even when they have the same spelling.
func (e *Emitter) Symbol(c spec.Callable, ref model.SymbolReference) string {
	var name string
	switch ref.Kind {
	case model.ParameterSymbol:
		for _, slot := range c.Params {
			if slot.ID == ref.ID {
				name = slot.Name
				break
			}
		}
	case model.ResultSymbol:
		for _, slot := range c.Results {
			if slot.ID == ref.ID {
				name = slot.Name
				break
			}
		}
	case model.ReceiverSymbol:
		name = c.Receiver
	case model.CancellationSymbol:
		if c.Cancellable {
			name = e.Profile.Cancellation.Name
		}
	case model.ErrorResultSymbol:
		if c.NamedResults && c.Fallible {
			name = e.Profile.ErrorResultName
		}
	case model.FieldSymbol:
		if typ, ok := e.Module.TypesByID[ref.Owner]; ok {
			for _, field := range typ.Fields {
				if field.ID == ref.ID {
					name = field.Name
					break
				}
			}
		}
	case model.LocalSymbol:
		name = ref.ID
	}
	if name == "" || !token.IsIdentifier(name) {
		e.Plugin.Error(fmt.Errorf("unresolved symbol %s in %s", ref.ID, c.ID))
		return "invalidSymbol"
	}
	return name
}
func (e *Emitter) Param(c spec.Callable, id string) string {
	return e.Symbol(c, model.SymbolReference{Kind: model.ParameterSymbol, ID: id})
}
func (e *Emitter) Result(c spec.Callable, id string) string {
	return e.Symbol(c, model.SymbolReference{Kind: model.ResultSymbol, ID: id})
}
func (e *Emitter) Receiver(c spec.Callable) string {
	return e.Symbol(c, model.SymbolReference{Kind: model.ReceiverSymbol})
}
func (e *Emitter) Context(c spec.Callable) string {
	return e.Symbol(c, model.SymbolReference{Kind: model.CancellationSymbol})
}
func (e *Emitter) ErrorResult(c spec.Callable) string {
	return e.Symbol(c, model.SymbolReference{Kind: model.ErrorResultSymbol})
}
func (e *Emitter) Field(owner, id string) string {
	return e.Symbol(spec.Callable{ID: owner}, model.SymbolReference{Kind: model.FieldSymbol, Owner: owner, ID: id})
}
func (e *Emitter) Local(name string) string {
	return e.Symbol(spec.Callable{ID: "local"}, model.SymbolReference{Kind: model.LocalSymbol, ID: name})
}

// Body emits supplied fragments verbatim. References must be resolved explicitly;
// emitted text is never scanned or renamed.
func (e *Emitter) Body(c spec.Callable, parts ...any) { e.G.P(parts...) }
