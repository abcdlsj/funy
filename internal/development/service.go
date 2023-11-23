package development

type Service struct {
	RunFn      func() error
	ShutdownFn func() error
	Name       string
}

func (s *Service) Run() error {
	return s.RunFn()
}

func (s *Service) Shutdown() error {
	return s.ShutdownFn()
}
