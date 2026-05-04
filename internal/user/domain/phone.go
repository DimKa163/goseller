package domain

import (
	"encoding/json"
	"errors"
	"regexp"
)

var phoneRe = regexp.MustCompile(`^[0-9]\d{7,14}$`)

var ErrInvalidPhone = errors.New("invalid phone format")

type Phone string

func (p Phone) Validate() error {
	if !phoneRe.MatchString(string(p)) {
		return ErrInvalidPhone
	}
	return nil
}

func (p Phone) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}

func (p *Phone) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	phone := Phone(value)
	if err := phone.Validate(); err != nil {
		return err
	}

	*p = phone
	return nil
}

func (p Phone) String() string {
	return string(p)
}
