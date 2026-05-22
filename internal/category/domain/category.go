package domain

import (
	"context"
	"time"

	"github.com/beevik/guid"
)

type CategoryID guid.Guid

func NewCategoryID(id string) (CategoryID, error) {
	guid, err := guid.ParseString(id)
	if err != nil {
		return CategoryID{}, newCategoryIDInvalidError(id, err)
	}
	return CategoryID(*guid), nil
}

func (c CategoryID) MarshalJSON() ([]byte, error) {
	uid := guid.Guid(c)
	return []byte(`"` + uid.String() + `"`), nil
}

func (c *CategoryID) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return newCategoryIDInvalidError(str, nil)
	}
	id, err := NewCategoryID(str[1 : len(str)-1])
	if err != nil {
		return err
	}
	*c = id
	return nil
}

func (c CategoryID) String() string {
	uid := guid.Guid(c)
	return uid.String()
}

type Category struct {
	ID          CategoryID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	Inactive    bool
}

func NewCategory(name, description string) *Category {
	uid := guid.New()
	return &Category{
		ID:          CategoryID(*uid),
		Name:        name,
		Description: description,
	}
}

type CategoryRepository interface {
	FindByID(ctx context.Context, id CategoryID) (*Category, error)
	Insert(ctx context.Context, category *Category) (CategoryID, error)
	Update(ctx context.Context, category *Category) (*Category, error)
	Deactivate(ctx context.Context, id CategoryID) error
}
