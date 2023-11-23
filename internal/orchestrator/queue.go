package orchestrator

func (q *Orchestrator) PushQueue(service *Service) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Waitings = append(q.Waitings, service)
}

func (q *Orchestrator) NextQueue() *Service {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.Waitings) == 0 {
		return nil
	}

	service := q.Waitings[0]

	q.Waitings = q.Waitings[1:]

	return service
}
