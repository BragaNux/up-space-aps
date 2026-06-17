package usecases

import (
	"errors"
	"strings"
)

// limite de ~4MB pra imagens em base64 (data URI), de olho em foto de perfil/aluno sem pesar demais
const maxImageDataLength = 5_500_000

var errImageTooLarge = errors.New("a imagem é muito grande (máximo de ~4MB)")

// usado pelos usecases de post/aluno/perfil pra rejeitar imagem base64 grande demais antes de salvar
func validateImage(value string) error {
	if strings.HasPrefix(value, "data:image/") && len(value) > maxImageDataLength {
		return errImageTooLarge
	}
	return nil
}
