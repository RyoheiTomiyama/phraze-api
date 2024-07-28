package card

import (
	"context"
	"math"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/domain"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/samber/lo"
)

type cardService struct{}

type ICardService interface {
	EvalSchedule(ctx context.Context, grade int, prevSchedule *domain.CardSchedule) (*domain.CardSchedule, error)
}

func NewService() ICardService {
	return &cardService{}
}

func (s *cardService) EvalSchedule(ctx context.Context, grade int, prevSchedule *domain.CardSchedule) (*domain.CardSchedule, error) {
	if prevSchedule == nil {
		return nil, errutil.New(errutil.CodeInternalError, "prevScheduleは必須")
	}
	factor := lo.Ternary(prevSchedule.Efactor == 0, 1.0, prevSchedule.Efactor)
	interval := lo.Ternary(prevSchedule.Interval == 0, 20, prevSchedule.Interval)

	interval, factor = calcEvaluation(grade, interval, factor)

	nextSchedule := prevSchedule
	nextSchedule.Efactor = factor
	nextSchedule.Interval = interval
	nextSchedule.ScheduleAt = time.Now().Add(time.Duration(interval) * time.Minute)

	return nextSchedule, nil
}

func calcEvaluation(grade, interval int, factor float64) (nextInterval int, nextFactor float64) {
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
