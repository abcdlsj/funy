package orchestrator

func (q *Orchestrator) GetAll() []*Instance {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.Actives
}

func (q *Orchestrator) Push(ins *Instance) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Waitings = append(q.Waitings, ins)
}

func (q *Orchestrator) Next() *Instance {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.Waitings) == 0 {
		return nil
	}

	ins := q.Waitings[0]

	q.Waitings = q.Waitings[1:]

	return ins
}

func (s *Orchestrator) Add(service Instance) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Actives = append(s.Actives, &service)
}

func (s *Orchestrator) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, v := range s.Actives {
		if v.InsName == name {
			s.Actives = append(s.Actives[:i], s.Actives[i+1:]...)
			return
		}
	}
}

func (s *Orchestrator) GetByName(name string) *Instance {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.Actives {
		if v.InsName == name {
			return v
		}
	}

	return nil
}
