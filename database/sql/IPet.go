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
	if (!pet.RazaID.Valid || pet.RazaID.Int32 <= 0) {
		return errors.New("razaID is required")
	} else if pet.Sexo == "" {
		return errors.New("sexo is required")
	} else if (!pet.FechaNacimiento.Valid || pet.FechaNacimiento.Time.String() == "") {
		return errors.New("fechaNacimiento is required")
	}

	// Validate sexo value
	if pet.Sexo != "H" && pet.Sexo != "M" {
		return errors.New("sexo debe ser H o M")
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
	// Check for not null fields
	if (!pet.RazaID.Valid || pet.RazaID.Int32 <= 0) {
		return errors.New("razaID is required")
	} else if pet.Sexo == "" {
		return errors.New("sexo is required")
	} else if (!pet.FechaNacimiento.Valid || pet.FechaNacimiento.Time.String() == "") {
		return errors.New("fechaNacimiento is required")
	} else if pet.ImgUrl.Valid && !isValidImageName(pet.ImgUrl.String) {
		return errors.New("formato de imagen inválido")
	}

	// Validate sexo value
	if (pet.Sexo != "H" && pet.Sexo != "M") {
		return errors.New("sexo debe ser H o M")
	}

	if hasSpecialCharacters(pet.Nombre.String) {
		return errors.New("nombre must not contain special characters")
	} else if !hasOnlyAlphanumericAndPunctuation(pet.Descripcion.String) {
		return errors.New("descripcion must not contain special characters")
	}
	
	// Check length
	if len(pet.Nombre.String) > 100 {
		return errors.New("nombre must be less than 50 characters")
	} else if len(pet.Descripcion.String) > 255 {
		return errors.New("descripcion must be less than 100 characters")
	} else if len(pet.Sexo) > 1 {
		return errors.New("sexo must be 1 character long")
	} else if len(pet.ImgUrl.String) > 50 {
		return errors.New("imgUrl must be less than 100 characters")
	}

	return nil
}

func (pet *Mascota) ValidateLogin() error {
	return errors.New("pets can't login")
}