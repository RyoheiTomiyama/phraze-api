package service

import "testing"

func TestCalcSchedule(t *testing.T) {
	t.Run("grade 5", func(t *testing.T) {
		s := cardService{}
		interval := float64(6)
		factor := float64(1)

		t.Log("input", interval, factor)

		for range 5 {

			interval, factor = s.CalcSchedule(5, interval, factor)
			t.Log("output", interval, factor)
		}
	})
	t.Run("grade 4", func(t *testing.T) {
		s := cardService{}
		interval := float64(6)
		factor := float64(1)

		t.Log("input", interval, factor)

		for range 5 {

			interval, factor = s.CalcSchedule(4, interval, factor)
			t.Log("output", interval, factor)
		}
	})
	t.Run("grade 3", func(t *testing.T) {
		s := cardService{}
		interval := float64(6)
		factor := float64(1)

		t.Log("input", interval, factor)

		for range 5 {

			interval, factor = s.CalcSchedule(3, interval, factor)
			t.Log("output", interval, factor)
		}
	})
}
