package db

import "strings"

func (insertVetParams *InsertNewVetParams) Trim() {
	insertVetParams.Nombre = strings.TrimSpace(insertVetParams.Nombre)
	insertVetParams.ApellidoP = strings.TrimSpace(insertVetParams.ApellidoP)
	insertVetParams.ApellidoM = strings.TrimSpace(insertVetParams.ApellidoM)
	insertVetParams.Email = strings.TrimSpace(insertVetParams.Email)
	insertVetParams.Telefono.String = strings.TrimSpace(insertVetParams.Telefono.String)
	insertVetParams.ImgUrl.String = strings.TrimSpace(insertVetParams.ImgUrl.String)
	insertVetParams.PasswordHash = strings.TrimSpace(insertVetParams.PasswordHash)
}