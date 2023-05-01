package db

import "strings"

func (vetPlace *Negocio) Trim() {
	vetPlace.NombreNegocio = strings.TrimSpace(vetPlace.NombreNegocio)
	vetPlace.Token = strings.TrimSpace(vetPlace.Token)
}

func (vetPlace *Negocio) DeleteBlankFields() {
	vetPlace.NombreNegocio = strings.ReplaceAll(vetPlace.NombreNegocio, " ", "")
	vetPlace.Token = strings.ReplaceAll(vetPlace.Token, " ", "")
}