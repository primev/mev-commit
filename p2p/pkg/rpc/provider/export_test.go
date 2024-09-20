package providerapi

func (s *Service) ActiveReceivers() int {
	return int(s.activeReceivers.Load())
}
