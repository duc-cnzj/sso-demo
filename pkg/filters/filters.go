package filters

import (
"github.com/jinzhu/gorm"
"reflect"
)

type GormScopeFunc = func(*gorm.DB) *gorm.DB

type Filterable interface {
	Apply() []GormScopeFunc
	GetFuncByName(string) func() GormScopeFunc
	All() interface{}
	Push(GormScopeFunc)
	Scopes() []GormScopeFunc
	ResetScopes()
}

func DefaultApply() func(f Filterable) []GormScopeFunc {
	return func(f Filterable) []GormScopeFunc {
		f.ResetScopes()
		refVaIn := reflect.ValueOf(f.All())
		for i := 0; i < refVaIn.Elem().NumField(); i++ {
			if !refVaIn.Elem().Field(i).IsZero() {
				if fn := f.GetFuncByName(refVaIn.Elem().Type().Field(i).Name); fn != nil {
					f.Push(fn())
				}
			}
		}

		var sos = make([]GormScopeFunc, len(f.Scopes()))
		copy(sos, f.Scopes())

		return sos
	}
}

