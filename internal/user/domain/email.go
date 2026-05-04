package domain

import (
	"encoding/json"
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z]{1}[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var ErrInvalidEmail = errors.New("invalid email format")

type Email string

func (e Email) Validate() error {
	if !emailRegex.MatchString(string(e)) {
		return ErrInvalidEmail
	}
	return nil
}

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

func (e *Email) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	email := Email(value)
	if err := email.Validate(); err != nil {
		return err
	}

	*e = email
	return nil
}

func (e Email) String() string {
	return string(e)
}
