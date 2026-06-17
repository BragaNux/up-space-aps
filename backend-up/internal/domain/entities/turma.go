package entities

import "time"

// Turma representa uma sala/classe da escola onde os alunos sao agrupados
type Turma struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
