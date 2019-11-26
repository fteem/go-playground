package sum

func Sum(nums []int) (sum int) {
	for _, num := range nums {
		sum += num
	}
	return sum
}

func SumAll(numsToSum ...[]int) (sums []int) {
	for _, numbers := range numsToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}

func SumAllTails(numsToSum ...[]int) (sums []int) {
	for _, numbers := range numsToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tails := numbers[1:]
			sums = append(sums, Sum(tails))
		}
	}

	return sums
}
