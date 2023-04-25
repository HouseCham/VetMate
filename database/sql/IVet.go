package db

import "strings"

func (vet *Veterinario) Trim() {
	vet.Nombre = strings.TrimSpace(vet.Nombre)
	vet.ApellidoP = strings.TrimSpace(vet.ApellidoP)
	vet.ApellidoM = strings.TrimSpace(vet.ApellidoM)
	vet.Email = strings.TrimSpace(vet.Email)
	vet.Telefono.String = strings.TrimSpace(vet.Telefono.String)
	vet.ImgUrl.String = strings.TrimSpace(vet.ImgUrl.String)
	vet.PasswordHash = strings.TrimSpace(vet.PasswordHash)
}