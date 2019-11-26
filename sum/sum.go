package sum

func Sum(nums []int) (sum int) {
	for _, num := range nums {
		sum += num
	}
	return sum
}
