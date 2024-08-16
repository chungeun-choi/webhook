package cert

import (
	"bytes"
	"encoding/pem"
	"github.com/chungeun-choi/webhook/errors"
)

func EncodeToPem(data []byte, typeInfo string) (*bytes.Buffer, error) {
	result := new(bytes.Buffer)
	if err := pem.Encode(result, &pem.Block{Type: typeInfo, Bytes: data}); err != nil {
		return nil, errors.FailedGenerateCertPEM(err)
	} else {
		return result, nil
	}
}
