package domain

import (
	"time"

	"github.com/beevik/guid"
)

type Category struct {
	ID          guid.Guid
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
}
