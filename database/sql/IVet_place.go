package db

import "strings"

func (vetPlace *Sucursale) Trim() {
	vetPlace.NombreSucursal = strings.TrimSpace(vetPlace.NombreSucursal)
	vetPlace.Token = strings.TrimSpace(vetPlace.Token)
}

func (vetPlace *Sucursale) DeleteBlankFields() {
	vetPlace.NombreSucursal = strings.ReplaceAll(vetPlace.NombreSucursal, " ", "")
	vetPlace.Token = strings.ReplaceAll(vetPlace.Token, " ", "")
}