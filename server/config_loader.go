package server

import (
	"github.com/rs/zerolog/log"
	"sso/config/env"
)

type ConfigLoader struct {
}

func (c *ConfigLoader) GetWeight() int {
	return 1
}

func (c *ConfigLoader) Load(s *Server) error {
	var (
		config *env.Config
		err    error
	)
	if config, err = ReadConfig(s.configPath); err != nil {
		log.Debug().Err(err).Msg("ConfigLoader Load")
		return err
	}
	s.config = config
	log.Info().Msg("Config loaded.")

	return nil
}
