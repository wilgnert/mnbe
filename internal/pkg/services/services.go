package services

const (
	nanosecond  = 1
	microsecond = 1000 * nanosecond
	milisecond  = 1000 * microsecond
	second      = 1000 * milisecond
	minute      = 60 * second
	hour        = 60 * minute
)

type ServiceError struct {
	msg string
}

func (s *ServiceError) Error() string {
	return s.msg
}

