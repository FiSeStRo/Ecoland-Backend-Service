package service

type HealthService interface {
	GetHealthStatus() string
}

type healthService struct {
}

func NewHealthService() *healthService {
	return &healthService{}
}

func (s *healthService) GetHealthStatus() string {
	return "Service is up and helthy"
}
