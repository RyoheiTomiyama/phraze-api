package card


type cardService struct{}

type ICardService interface {
	CalcSchedule(grade int, interval int, factor float64) (nextInterval int, nextFactor float64)
}

func NewCardService() ICardService {
	return &cardService{}
}

func (s *cardService) CalcSchedule(grade int, interval int, factor float64) (nextInterval int, nextFactor float64) {
	factor += (0.2 - float64(5-grade)*(0.08+float64(5-grade)*0.02)) / (factor * factor * factor)
	factor = round(max(1.0, factor), 2)

	interval = int(float64(interval) * float64(grade-2) * factor)

	return interval, factor
}

// pos桁数で四捨五入
func round(num float64, pos int) float64 {
	shift := math.Pow10(pos)             // 小数の位置をずらすためのシフト値を算出
	shiftedNum := num * shift            // 四捨五入したい桁を小数第一位にずらす
	roundedNum := math.Round(shiftedNum) // 小数第一位を四捨五入する
	result := roundedNum / shift         // 小数の位置を元に戻す
	return result
}
