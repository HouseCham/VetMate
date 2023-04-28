package interfaces

type IDatabaseValidation interface {
	ValidateNewRegister() error
	ValidateUpdate() error
	ValidateLogin() error
}

type INewInsertParams interface {
	Trim()
	DeleteBlankFields()
}