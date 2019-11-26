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
