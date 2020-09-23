package server

import "sso/config/env"

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
		return err
	}
	s.config = config

	return nil
}
