package sum

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("Got '%d', want '%d'", got, want)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("Got '%d', want '%d'", got, want)
		}
	})
}

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2}, []int{1, 2})
	want := []int{3, 3}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got '%d', want '%d'", got, want)
	}
}

func TestSumAllTails(t *testing.T) {
	checkSums := func(t *testing.T, got, want []int) {
		t.Helper()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got '%d', want '%d'", got, want)
		}
	}

	t.Run("sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{1, 2})
		want := []int{2, 2}
		checkSums(t, got, want)
	})

	t.Run("sums of empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{})
		want := []int{0, 0}
		checkSums(t, got, want)
	})
}
