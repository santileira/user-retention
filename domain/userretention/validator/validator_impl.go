package validator

import "fmt"

type UserRetentionValidatorImpl struct{}

func NewUserRetentionValidatorImpl() *UserRetentionValidatorImpl {
	return &UserRetentionValidatorImpl{}
}

func (u *UserRetentionValidatorImpl) ValidateInput(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path can't be empty")
	}

	return nil
}
