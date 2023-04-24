package db

import "strings"

func (insertNewUserParams *InsertNewUserParams) Trim(){
	insertNewUserParams.Nombre = strings.TrimSpace(insertNewUserParams.Nombre)
	insertNewUserParams.ApellidoP = strings.TrimSpace(insertNewUserParams.ApellidoP)
	insertNewUserParams.ApellidoM = strings.TrimSpace(insertNewUserParams.ApellidoM)
	insertNewUserParams.Email = strings.TrimSpace(insertNewUserParams.Email)
	insertNewUserParams.Telefono.String = strings.TrimSpace(insertNewUserParams.Telefono.String)
}