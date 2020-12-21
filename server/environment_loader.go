package server

import "github.com/rs/zerolog/log"

type EnvironmentLoader struct {
}

func (e *EnvironmentLoader) Load(s *Server) error {
	if s.env.IsProduction() {
		s.ProductionMode()
	} else if s.env.IsDebugging() {
		s.DebugMode()
	}

	log.Info().Msg("Environment loaded.")

	return nil
}

func (e *EnvironmentLoader) GetWeight() int {
	return 8
}
