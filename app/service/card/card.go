package service

type cardService struct{}

type ICardService interface{}

func NewCardService() ICardService {
	return &cardService{}
}

func (s *cardService) CalcSchedule(grade int, interval, factor float64) (nextInterval, nextFactor float64) {
	factor += 0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02)
	factor = max(1.3, factor)

	interval = interval * factor

	return interval, factor
}
