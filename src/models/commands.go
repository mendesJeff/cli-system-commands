package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type (
	StringInterfaceMap map[string]interface{}
	Command            struct {
		ID         string `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		Parameters string `json:"parameters,omitempty"`
		Olt_id     uint64 `json:"olt_id,omitempty"`
		Client_id  uint64 `json:"client_id,omitempty"`
		Olt_name   string `json:"olt_name,omitempty"`
		Onu_serial string `json:"onu_serial,omitempty"`
		Error      string `json:"error,omitempty"`
		/* Parameters    StringInterfaceMap `json:"parameters,omitempty"` */
		Last_update   time.Time `json:"last_update,omitempty"`
		Creation_date time.Time `json:"creation_date,omitempty"`
	}
)

type Scanner interface {
	Scan(src interface{}) error
}

func (m *StringInterfaceMap) Scan(src interface{}) error {
	var source []byte
	_m := make(map[string]interface{})

	switch src.(type) {
	case []uint8:
		source = []byte(src.([]uint8))
	case nil:
		return nil
	default:
		return errors.New("incompatible type for StringInterfaceMap")
	}
	err := json.Unmarshal(source, &_m)
	if err != nil {
		return err
	}
	*m = StringInterfaceMap(_m)
	return nil
}

func (m StringInterfaceMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return driver.Value([]byte(j)), nil
}

func (command *Command) Prepare() error {
	if error := command.validate(); error != nil {
		return error
	}

	command.formate()
	return nil
}

func (command *Command) validate() error {
	switch {
	case command.Name == "":
		return errors.New("the field 'name' is mandatory")
	case command.Parameters == "":
		return errors.New("the field 'parameters' is mandatory")

	}
	return nil
}

func (command *Command) formate() {
	command.Name = strings.ReplaceAll(command.Name, " ", "")
	command.Name = strings.TrimSpace(command.Name)
	command.Parameters = strings.TrimSpace(command.Parameters)

}
