package env

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/jinzhu/gorm"
)

type Env struct {
	translator *ut.UniversalTranslator
	db         *gorm.DB
}

type Operator func(*Env)

func NewEnv(db *gorm.DB, envOperators ...Operator) *Env {
	env := &Env{
		db: db,
	}
	for _, op := range envOperators {
		op(env)
	}

	return env
}

func WithUniversalTranslator(t *ut.UniversalTranslator) func(env *Env) {
	return func(env *Env) {
		env.translator = t
	}
}

func (e *Env) GetUniversalTranslator() *ut.UniversalTranslator {
	return e.translator
}

func (e *Env) GetDB() *gorm.DB {
	return e.db
}

