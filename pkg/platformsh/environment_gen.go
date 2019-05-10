// This file is generated - do not edit!

package platformsh

import (
	"encoding/base64"
	"encoding/json"
)

type EnvironmentAPI interface {
	Application() (*Application, error)
	ApplicationName() (string, error)
	AppName() (string, error)
	AppCommand() (string, error)
	ApplicationCommand() (string, error)
	AppDir() (string, error)
	Branch() (string, error)
	Dir() (string, error)
	DocumentRoot() (string, error)
	Environment() (string, error)
	Port() (string, error)
	Project() (string, error)
	ProjectEntropy() (string, error)
	Relationships() (Relationships, error)
	Routes() (Routes, error)
	SMTPHost() (string, error)
	Socket() (string, error)
	TreeID() (string, error)
	Variables() (JSONObject, error)
	Vars() (JSONObject, error)
	XClientCert() (string, error)
	XClientDN() (string, error)
	XClientIP() (string, error)
	XClientSSL() (string, error)
	XClientVerify() (string, error)
}

func (e *environment) Application() (*Application, error) {
	name := e.Prefix() + "APPLICATION"
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

func (e *environment) ApplicationName() (string, error) {
	name := e.Prefix() + "APPLICATION_NAME"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) AppName() (string, error) {
	name := e.Prefix() + "APPLICATION_NAME"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) AppCommand() (string, error) {
	name := e.Prefix() + "APP_COMMAND"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) ApplicationCommand() (string, error) {
	name := e.Prefix() + "APP_COMMAND"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) AppDir() (string, error) {
	name := e.Prefix() + "APP_DIR"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Branch() (string, error) {
	name := e.Prefix() + "BRANCH"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Dir() (string, error) {
	name := e.Prefix() + "DIR"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) DocumentRoot() (string, error) {
	name := e.Prefix() + "DOCUMENT_ROOT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Environment() (string, error) {
	name := e.Prefix() + "ENVIRONMENT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Port() (string, error) {
	name := "PORT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Project() (string, error) {
	name := e.Prefix() + "PROJECT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) ProjectEntropy() (string, error) {
	name := e.Prefix() + "PROJECT_ENTROPY"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Relationships() (Relationships, error) {
	name := e.Prefix() + "RELATIONSHIPS"
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

func (e *environment) Routes() (Routes, error) {
	name := e.Prefix() + "ROUTES"
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

func (e *environment) SMTPHost() (string, error) {
	name := e.Prefix() + "SMTP_HOST"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Socket() (string, error) {
	name := "SOCKET"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) TreeID() (string, error) {
	name := e.Prefix() + "TREE_ID"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) Variables() (JSONObject, error) {
	name := e.Prefix() + "VARIABLES"
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

func (e *environment) Vars() (JSONObject, error) {
	name := e.Prefix() + "VARIABLES"
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

func (e *environment) XClientCert() (string, error) {
	name := "X_CLIENT_CERT"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) XClientDN() (string, error) {
	name := "X_CLIENT_DN"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) XClientIP() (string, error) {
	name := "X_CLIENT_IP"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) XClientSSL() (string, error) {
	name := "X_CLIENT_SSL"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}

func (e *environment) XClientVerify() (string, error) {
	name := "X_CLIENT_VERIFY"
	value, ok := e.lookup(name)
	if !ok {
		return "", missingEnvironment(name)
	}

	return value, nil
}
