package orchestrator

func (s *Orchestrator) AddService(service Service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Services = append(s.Services, &service)
}

func (s *Orchestrator) RemoveService(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.Services {
		if v.Instance.Name == name {
			s.Services = append(s.Services[:i], s.Services[i+1:]...)
			return
		}
	}
}

func (s *Orchestrator) GetService(name string) *Service {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.Services {
		if v.Instance.Name == name {
			return v
		}
	}

	return nil
}
