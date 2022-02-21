package calculator

type UserRetentionCalculator interface {
	Calculate(rawUsersActivity [][]string) (map[int][]int, error)
}
