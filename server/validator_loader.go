package server

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"regexp"
)

type ValidatorLoader struct {
}

func (v *ValidatorLoader) Load(s *Server) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
			slug, ok := fl.Field().Interface().(string)
			if ok {
				regex := regexp.MustCompile("^[a-zA-Z_-]+$")
				match := regex.Match([]byte(slug))
				if match {
					return true
				}
			}

			return false
		})
	}

	log.Info().Msg("Validator loaded.")

	return nil
}

func (v *ValidatorLoader) GetWeight() int {
	return 10
}
