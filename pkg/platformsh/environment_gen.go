// This file is generated - do not edit!

package platformsh

import (
	"encoding/base64"
	"encoding/json"
)

func (e Environment) Application() (*Application, error) {
	name := e.Prefix + "APPLICATION"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := Application{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return &obj, nil
}

func (e Environment) ApplicationName() (string, error) {
	name := e.Prefix + "APPLICATION_NAME"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) AppName() (string, error) {
	name := e.Prefix + "APPLICATION_NAME"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) AppCommand() (string, error) {
	name := e.Prefix + "APP_COMMAND"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) ApplicationCommand() (string, error) {
	name := e.Prefix + "APP_COMMAND"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) AppDir() (string, error) {
	name := e.Prefix + "APP_DIR"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Branch() (string, error) {
	name := e.Prefix + "BRANCH"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Dir() (string, error) {
	name := e.Prefix + "DIR"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) DocumentRoot() (string, error) {
	name := e.Prefix + "DOCUMENT_ROOT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Environment() (string, error) {
	name := e.Prefix + "ENVIRONMENT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Port() (string, error) {
	name := "PORT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Project() (string, error) {
	name := e.Prefix + "PROJECT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) ProjectEntropy() (string, error) {
	name := e.Prefix + "PROJECT_ENTROPY"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Relationships() (Relationships, error) {
	name := e.Prefix + "RELATIONSHIPS"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := Relationships{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return obj, nil
}

func (e Environment) Routes() (Routes, error) {
	name := e.Prefix + "ROUTES"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := Routes{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return obj, nil
}

func (e Environment) SMTPHost() (string, error) {
	name := e.Prefix + "SMTP_HOST"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Socket() (string, error) {
	name := "SOCKET"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) TreeID() (string, error) {
	name := e.Prefix + "TREE_ID"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e Environment) Variables() (JSONObject, error) {
	name := e.Prefix + "VARIABLES"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := JSONObject{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return obj, nil
}

func (e Environment) Vars() (JSONObject, error) {
	name := e.Prefix + "VARIABLES"
	value, ok := e.lookup(name)
	if !ok {
		return nil, missingEnvironment(name)
	}

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	obj := JSONObject{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, MissingEnvironment{name, err}
	}

	return obj, nil
}
