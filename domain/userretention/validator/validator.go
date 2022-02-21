package validator

type UserRetentionValidator interface {
	ValidateInput(filePath string) error
}
