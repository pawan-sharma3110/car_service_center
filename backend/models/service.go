package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ServiceID uuid.UUID `json:"service_id"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
	AddedOn   time.Time `json:"added_on"`
}
