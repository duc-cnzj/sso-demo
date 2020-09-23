package server

import "sso/config/env"

type EnvLoader struct {
}

func (e *EnvLoader) GetWeight() int {
	return 5
}

func (e *EnvLoader) Load(s *Server) error {
	s.env = env.NewEnv(
		s.config,
		s.db,
		s.session,
		s.redis,
		env.WithUniversalTranslator(s.LoadTranslators()),
		env.WithRootDir(s.rootPath),
	)

	return nil
}
