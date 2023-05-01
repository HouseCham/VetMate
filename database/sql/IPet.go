package db

import (
	"errors"
	"strings"
)

func (pet *Mascota) Trim() {
	pet.Nombre.String = strings.TrimSpace(pet.Nombre.String)
	pet.Descripcion.String = strings.TrimSpace(pet.Descripcion.String)
	pet.Sexo = strings.TrimSpace(pet.Sexo)
	pet.Token = strings.TrimSpace(pet.Token)
	pet.ImgUrl.String = strings.TrimSpace(pet.ImgUrl.String)
}

func (pet *Mascota) DeleteBlankFields() {
	pet.Nombre.String = strings.ReplaceAll(pet.Nombre.String, " ", "")
	pet.Descripcion.String = strings.ReplaceAll(pet.Descripcion.String, " ", "")
	pet.Sexo = strings.ReplaceAll(pet.Sexo, " ", "")
	pet.Token = strings.ReplaceAll(pet.Token, " ", "")
	pet.ImgUrl.String = strings.ReplaceAll(pet.ImgUrl.String, " ", "")
}

func (pet *Mascota) ValidateNewRegister() error {
	// Check for not null fields
	if !pet.RazaID.Valid {
		return errors.New("razaID is required")
	} else if pet.Sexo == "" {
		return errors.New("sexo is required")
	} 

	// Check length
	if len(pet.Nombre.String) > 100 {
		return errors.New("nombre must be less than 50 characters")
	} else if len(pet.Descripcion.String) > 255 {
		return errors.New("descripcion must be less than 100 characters")
	} else if len(pet.Sexo) > 1 {
		return errors.New("sexo must be 1 character long")
	} else if len(pet.Token) > 10 {
		return errors.New("token must be less than 10 characters")
	} else if len(pet.ImgUrl.String) > 50 {
		return errors.New("imgUrl must be less than 100 characters")
	}

	return nil
}

func (pet *Mascota) ValidateUpdate() error {
	return errors.New("implement me")
}

func (pet *Mascota) ValidateLogin() error {
	return errors.New("pets can't login")
}