package cert

import (
	"bytes"
)

type Info struct {
	Org      []string
	DNSNames []string
	Commons  string
	KeyPath  string
	CAPath   string
	CaType   int
}

type Certificate struct {
	CA         *bytes.Buffer
	Cert       *bytes.Buffer
	PrivateKey *bytes.Buffer
}

var CertTypes []int = []int{SelfSigned, Private, Public}

// CA cert type
const (
	SelfSigned = iota
	Private
	Public
)
