package interfaces

type IDatabase interface {
	InsertNew() error
	Update() error
	Delete() error
	FindAll() error
	FindById() error
}

type INewInsertParams interface {
	Trim()
}