package pki

type Bundle struct {
	Name string
	Cert Certificate
	Key  PrivateKey
}
