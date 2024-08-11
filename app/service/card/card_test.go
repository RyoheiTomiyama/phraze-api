package card

import (
	"testing"
	"time"
)

func TestCalcEvaluation(t *testing.T) {
	t.Run("grade 5", func(t *testing.T) {
		interval := 20
		factor := float64(1)

		t.Log("input", interval, factor)

		for range 5 {
			interval, factor = calcEvaluation(5, interval, factor)
			t.Log("output", interval, factor)
			t.Log(
				"duration", (time.Duration(interval) * time.Minute).String(),
				"next", time.Now().Add(time.Duration(interval)*time.Minute).Format(time.RFC3339Nano),
			)
		}
	})
	t.Run("grade 4", func(t *testing.T) {
		interval := 20
		factor := float64(1)

		t.Log("input", interval, factor)

		for range 5 {
			interval, factor = calcEvaluation(4, interval, factor)
			t.Log("output", interval, factor)
			t.Log(
				"duration", (time.Duration(interval) * time.Minute).String(),
				"next", time.Now().Add(time.Duration(interval)*time.Minute).Format(time.RFC3339Nano),
			)
		}
	})
	t.Run("grade 3", func(t *testing.T) {
		interval := 20
		factor := float64(1)

		t.Log("input", interval, factor)

		for range 5 {
			interval, factor = calcEvaluation(3, interval, factor)
			t.Log("output", interval, factor)
			t.Log(
				"duration", (time.Duration(interval) * time.Minute).String(),
				"next", time.Now().Add(time.Duration(interval)*time.Minute).Format(time.RFC3339Nano),
			)
		}
	})
}
