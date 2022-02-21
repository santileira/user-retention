package calculator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateUserRetention(t *testing.T) {
	cases := []struct {
		name               string
		userRetentionInput map[int][]int
		userActivityInput  []int
		expectedOutput     map[int][]int
	}{
		{
			name: "Should update user retention when user retention is empty and user has activity in the first day",
			userRetentionInput: func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				return userRetention
			}(),
			userActivityInput: []int{1},
			expectedOutput: func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[0][0] = 1
				return userRetention
			}(),
		},
		{
			name: "Should update user retention when user retention is empty and user has activity in the last day",
			userRetentionInput: func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				return userRetention
			}(),
			userActivityInput: []int{14},
			expectedOutput: func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[13][0] = 1
				return userRetention
			}(),
		},
		{
			name: "Should update user retention when user has activity between 2 and 5 days",
			userRetentionInput: func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[0][5] = 200
				userRetention[1][3] = 199
				return userRetention
			}(),
			userActivityInput: []int{2, 3, 4, 5},
			expectedOutput: func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[0][5] = 200
				userRetention[1][3] = 200
				return userRetention
			}(),
		},
	}

	userRetentionCalculator := NewUserRetentionCalculatorImpl()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// Operation
			userRetentionCalculator.updateUserRetention(c.userRetentionInput, c.userActivityInput)

			// Validation
			assert.EqualValues(t, c.expectedOutput, c.userRetentionInput)
		})
	}
}

func TestGetDayFromRawUserActivity(t *testing.T) {
	cases := []struct {
		name           string
		input          []string
		expectedOutput int
		expectedError  error
	}{
		{
			"Should return an error when value in the first position is not an int",
			[]string{"invalid value", "userid"},
			0,
			fmt.Errorf("error converting timestamp invalid value, err: strconv.Atoi: parsing \"invalid value\": invalid syntax"),
		},
		{
			"Should return an error when value in the first position is not an int",
			[]string{"1609459200", "userid"},
			1,
			nil,
		},
	}

	userRetentionCalculator := NewUserRetentionCalculatorImpl()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// Operation
			dayInUTC, err := userRetentionCalculator.getDayFromRawUserActivity(c.input)

			// Validation
			assert.EqualValues(t, c.expectedOutput, dayInUTC)
			assert.EqualValues(t, c.expectedError, err)
		})
	}
}

func TestCalculate(t *testing.T) {
	cases := []struct {
		name           string
		input          [][]string
		expectedOutput map[int][]int
		expectedError  error
	}{
		{
			"Should return default value when raw users activity is nil",
			nil,
			func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				return userRetention
			}(),
			nil,
		},
		{
			"Should return default value when raw users activity is empty",
			[][]string{},
			func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				return userRetention
			}(),
			nil,
		},
		{
			"Should return error when raw users activity has wrong data",
			[][]string{{"timestamp1", "userid1"}},
			nil,
			fmt.Errorf("error converting timestamp timestamp1, err: strconv.Atoi: parsing \"timestamp1\": invalid syntax"),
		},
		{
			"Should return user retention when there's just one user with activity on the first day",
			[][]string{{"1609459200", "1"}},
			func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[0][0] = 1
				return userRetention
			}(),
			nil,
		},
		{
			"Should return user retention when there's just one user with activity on the last day",
			[][]string{{"1609804800", "5"}},
			func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[4][0] = 1
				return userRetention
			}(),
			nil,
		},
		{
			"Should return user retention when there's just one user with two activities on the same day",
			[][]string{{"1609459200", "1", "1609459260", "1"}},
			func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[0][0] = 1
				return userRetention
			}(),
			nil,
		},
		{
			"Should return user retention when there's multiple users",
			[][]string{
				{"1609459200", "1"},
				{"1609459200", "2"},
				{"1609459200", "3"},
				{"1609459200", "4"},
				{"1609459260", "1"},
				{"1609545600", "1"},
				{"1609545600", "3"},
				{"1609632000", "1"},
				{"1609632000", "2"},
				{"1609632000", "3"},
				{"1609718400", "1"},
				{"1609718400", "2"},
				{"1609804800", "1"},
				{"1609804800", "5"},
			},
			func() map[int][]int {
				userRetention := make(map[int][]int, 14)
				for i := 0; i < 14; i++ {
					userRetention[i] = make([]int, 14)
				}
				userRetention[0][0] = 2
				userRetention[0][2] = 1
				userRetention[0][4] = 1
				userRetention[2][1] = 1
				userRetention[4][0] = 1
				return userRetention
			}(),
			nil,
		},
	}

	userRetentionCalculator := NewUserRetentionCalculatorImpl()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// Operation
			userRetention, err := userRetentionCalculator.Calculate(c.input)

			// Validation
			assert.EqualValues(t, c.expectedOutput, userRetention)
			assert.EqualValues(t, c.expectedError, err)
		})
	}
}
