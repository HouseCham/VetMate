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

func (updateVetParams *UpdateVetParams) Trim() {
	updateVetParams.Nombre = strings.TrimSpace(updateVetParams.Nombre)
	updateVetParams.ApellidoP = strings.TrimSpace(updateVetParams.ApellidoP)
	updateVetParams.ApellidoM = strings.TrimSpace(updateVetParams.ApellidoM)
	updateVetParams.Telefono.String = strings.TrimSpace(updateVetParams.Telefono.String)
	updateVetParams.ImgUrl.String = strings.TrimSpace(updateVetParams.ImgUrl.String)
}