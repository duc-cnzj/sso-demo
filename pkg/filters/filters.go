package filters

import (
	"github.com/jinzhu/gorm"
)

type GormScopeFunc = func(*gorm.DB) *gorm.DB

type Filterable interface {
	Apply() []GormScopeFunc
	GetFuncByName(string) func() GormScopeFunc
	All() []string
	Push(GormScopeFunc)
	Scopes() []GormScopeFunc
	ResetScopes()
}

func DefaultApply() func(f Filterable) []GormScopeFunc {
	return func(f Filterable) []GormScopeFunc {
		f.ResetScopes()
		for _, key := range f.All() {
			if fn := f.GetFuncByName(key); fn != nil {
				f.Push(fn())
			}
		}

		var sos = make([]GormScopeFunc, len(f.Scopes()))
		copy(sos, f.Scopes())

		return sos
	}
}

