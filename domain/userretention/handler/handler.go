package handler

import (
	"fmt"
	"github.com/santileira/user-retention/domain/userretention/calculator"
	"github.com/santileira/user-retention/domain/userretention/validator"
	"github.com/santileira/user-retention/filereader"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type UserRetentionHandler struct {
	validator  validator.UserRetentionValidator
	fileReader filereader.FileReader
	calculator calculator.UserRetentionCalculator
}

func NewUserRetentionHandler(validator validator.UserRetentionValidator, fileReader filereader.FileReader, calculator calculator.UserRetentionCalculator) *UserRetentionHandler {
	return &UserRetentionHandler{
		validator:  validator,
		fileReader: fileReader,
		calculator: calculator,
	}
}

func (u *UserRetentionHandler) HandleRequest(filePath string) (string, error) {
	result := ""
	if err := u.validator.ValidateInput(filePath); err != nil {
		logrus.Errorf("Error validating the input, err: %s", err.Error())
		return result, err
	}

	rawUsersActivity, err := u.fileReader.OpenFile(filePath)
	if err != nil {
		logrus.Errorf("Error opening the file, err: %s", err.Error())
		return result, err
	}

	usersRetention, err := u.calculator.Calculate(rawUsersActivity)
	if err != nil {
		logrus.Errorf("Error calculating users retention, err: %s", err.Error())
		return result, err
	}

	for i := 0; i < len(usersRetention); i++ {
		values := usersRetention[i]
		result += fmt.Sprintf("%d,%s\n", i+1, u.convertSliceOfIntsToString(values))
	}
	return result, nil
}

var res = 0

func (u *UserRetentionHandler) convertSliceOfIntsToString(ints []int) string {
	strs := make([]string, len(ints))
	for index, value := range ints {
		res += value
		strs[index] = strconv.Itoa(value)
	}

	return strings.Join(strs, ",")
}
