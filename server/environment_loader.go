package server

type EnvironmentLoader struct {
}

func (e *EnvironmentLoader) Load(s *Server) error {
	if s.env.IsProduction() {
		s.ProductionMode()
	} else if s.env.IsDebugging() {
		s.DebugMode()
	}

	return nil
}

func (e *EnvironmentLoader) GetWeight() int {
	return 8
}
