package handler

import (
	"fmt"
	"github.com/santileira/user-retention/domain/userretention/calculator"
	calculatormocks "github.com/santileira/user-retention/domain/userretention/calculator/mocks"
	"github.com/santileira/user-retention/domain/userretention/validator"
	validatormocks "github.com/santileira/user-retention/domain/userretention/validator/mocks"
	"github.com/santileira/user-retention/filereader"
	filereadermocks "github.com/santileira/user-retention/filereader/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestConvertSliceOfIntsToString(t *testing.T) {
	cases := []struct {
		name           string
		input          []int
		expectedOutput string
	}{
		{
			"Should return empty string when input is nil",
			nil,
			"",
		},
		{
			"Should return empty when input is empty",
			[]int{},
			"",
		},
		{
			"Should return the value when input has one value",
			[]int{0},
			"0",
		},
		{
			"Should return values separated with comma when input has more than one value",
			[]int{0, 1, 2, 3},
			"0,1,2,3",
		},
	}

	userRetentionHandler := NewUserRetentionHandler(nil, nil, nil)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// Operation
			output := userRetentionHandler.convertSliceOfIntsToString(c.input)

			// Validation
			assert.EqualValues(t, c.expectedOutput, output)
		})
	}
}

func TestHandleRequest(t *testing.T) {
	cases := []struct {
		name                        string
		input                       string
		mockUserRetentionValidator  *validatormocks.UserRetentionValidator
		mockFileReader              *filereadermocks.FileReader
		mockUserRetentionCalculator *calculatormocks.UserRetentionCalculator
		expectedOutput              string
		expectedError               error
	}{
		{
			"Should return an error when validator returns an error",
			"filePath",
			func() *validatormocks.UserRetentionValidator {
				mock := &validatormocks.UserRetentionValidator{}
				mock.On("ValidateInput", "filePath").
					Return(fmt.Errorf("error validating input"))
				return mock
			}(),
			nil,
			nil,
			"",
			fmt.Errorf("error validating input"),
		},
		{
			"Should return an error when file reader returns an error",
			"filePath",
			func() *validatormocks.UserRetentionValidator {
				mock := &validatormocks.UserRetentionValidator{}
				mock.On("ValidateInput", "filePath").
					Return(nil)
				return mock
			}(),
			func() *filereadermocks.FileReader {
				mock := &filereadermocks.FileReader{}
				mock.On("OpenFile", "filePath").
					Return(nil, fmt.Errorf("error opening the file"))
				return mock
			}(),
			nil,
			"",
			fmt.Errorf("error opening the file"),
		},
		{
			"Should return an error when calculator returns an error",
			"filePath",
			func() *validatormocks.UserRetentionValidator {
				mock := &validatormocks.UserRetentionValidator{}
				mock.On("ValidateInput", "filePath").
					Return(nil)
				return mock
			}(),
			func() *filereadermocks.FileReader {
				mock := &filereadermocks.FileReader{}
				mock.On("OpenFile", "filePath").
					Return([][]string{{"1609459200", "1"}, {"1609545600", "3"}}, nil)
				return mock
			}(),
			func() *calculatormocks.UserRetentionCalculator {
				mock := &calculatormocks.UserRetentionCalculator{}
				mock.On("Calculate", [][]string{{"1609459200", "1"}, {"1609545600", "3"}}).
					Return(nil, fmt.Errorf("error calculating user retention"))
				return mock
			}(),
			"",
			fmt.Errorf("error calculating user retention"),
		},
		{
			"Should return a nil error",
			"filePath",
			func() *validatormocks.UserRetentionValidator {
				mock := &validatormocks.UserRetentionValidator{}
				mock.On("ValidateInput", "filePath").
					Return(nil)
				return mock
			}(),
			func() *filereadermocks.FileReader {
				mock := &filereadermocks.FileReader{}
				mock.On("OpenFile", "filePath").
					Return([][]string{{"1609459200", "1"}, {"1609545600", "3"}}, nil)
				return mock
			}(),
			func() *calculatormocks.UserRetentionCalculator {
				userRetention := make(map[int][]int, 2)
				for i := 0; i < 2; i++ {
					userRetention[i] = make([]int, 2)
				}
				userRetention[0][0] = 1
				userRetention[1][0] = 1
				mock := &calculatormocks.UserRetentionCalculator{}
				mock.On("Calculate", [][]string{{"1609459200", "1"}, {"1609545600", "3"}}).
					Return(userRetention, nil)
				return mock
			}(),
			"1,1,0\n2,1,0\n",
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			userRetentionHandler := NewUserRetentionHandler(c.mockUserRetentionValidator,
				c.mockFileReader, c.mockUserRetentionCalculator)

			// Operation
			output, err := userRetentionHandler.HandleRequest(c.input)

			// Validation
			assert.EqualValues(t, c.expectedOutput, output)
			assert.EqualValues(t, c.expectedError, err)

			if c.mockUserRetentionValidator != nil {
				c.mockUserRetentionValidator.AssertNumberOfCalls(t, "ValidateInput", 1)
			}

			if c.mockFileReader != nil {
				c.mockFileReader.AssertNumberOfCalls(t, "OpenFile", 1)
			}

			if c.mockUserRetentionCalculator != nil {
				c.mockUserRetentionCalculator.AssertNumberOfCalls(t, "Calculate", 1)
			}
		})
	}
}

func TestHandleRequestIntegration(t *testing.T) {
	// Initialization
	file, _ := ioutil.TempFile("", "integration_test.csv")
	_, _ = file.Write([]byte("1609459200,1\n" +
		"1609459200,2\n" +
		"1609459200,3\n" +
		"1609459200,4\n" +
		"1609459260,1\n" +
		"1609545600,1\n" +
		"1609545600,3\n" +
		"1609632000,1\n" +
		"1609632000,2\n" +
		"1609632000,3\n" +
		"1609718400,1\n" +
		"1609718400,2\n" +
		"1609804800,1\n" +
		"1609804800,5"),
	)

	userRetentionValidator := validator.NewUserRetentionValidatorImpl()
	fileReader := filereader.NewFileReaderImpl()
	userRetentionCalculator := calculator.NewUserRetentionCalculatorImpl()
	userRetentionHandler := NewUserRetentionHandler(userRetentionValidator, fileReader, userRetentionCalculator)

	// Operation
	output, err := userRetentionHandler.HandleRequest(file.Name())

	// Validation
	assert.EqualValues(t, "1,2,0,1,0,1,0,0,0,0,0,0,0,0,0\n"+
		"2,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"3,0,1,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"4,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"5,1,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"6,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"7,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"8,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"9,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"10,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"11,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"12,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"13,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n"+
		"14,0,0,0,0,0,0,0,0,0,0,0,0,0,0\n", output)
	assert.Nil(t, err)
}
