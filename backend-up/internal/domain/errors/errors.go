package errors

import "errors"

// erros de dominio reutilizados pelos usecases pra nao ficar repetindo errors.New espalhado
var (
	ErrInvalidPresenceStatus = errors.New("status de presença inválido") // usado quando tentam marcar presenca com um status que nao existe
)
