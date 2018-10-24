package apis

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type VarsCertLoader struct {
	vars Variables
}

func NewVarsCertLoader(vars Variables) VarsCertLoader {
	return VarsCertLoader{vars}
}

func (l VarsCertLoader) LoadCerts(name string) (*x509.Certificate, *rsa.PrivateKey, error) {
	val, found, err := l.vars.Get(VariableDefinition{Name: name})
	if err != nil {
		return nil, nil, err
	} else if !found {
		return nil, nil, fmt.Errorf("Expected to find variable '%s' with a certificate", name)
	}

	// Convert to YAML for easier struct parsing
	valBytes, err := yaml.Marshal(val)
	if err != nil {
		return nil, nil, errors.New(err.Error() + fmt.Sprintf("Expected variable '%s' to be serializable", name))
	}

	type CertVal struct {
		Certificate string
		PrivateKey  string `yaml:"private_key"`
	}

	var certVal CertVal

	err = yaml.Unmarshal(valBytes, &certVal)
	if err != nil {
		return nil, nil, errors.New(err.Error() + fmt.Sprintf("Expected variable '%s' to be deserializable", name))
	}

	crt, err := l.parseCertificate(certVal.Certificate)
	if err != nil {
		return nil, nil, err
	}

	key, err := l.parsePrivateKey(certVal.PrivateKey)
	if err != nil {
		return nil, nil, err
	}

	return crt, key, nil
}

func (VarsCertLoader) parseCertificate(data string) (*x509.Certificate, error) {
	cpb, _ := pem.Decode([]byte(data))
	if cpb == nil {
		return nil, errors.New("Certificate did not contain PEM formatted block")
	}

	crt, err := x509.ParseCertificate(cpb.Bytes)
	if err != nil {
		return nil, errors.New(err.Error() + "Parsing certificate")
	}

	return crt, nil
}

func (VarsCertLoader) parsePrivateKey(data string) (*rsa.PrivateKey, error) {
	kpb, _ := pem.Decode([]byte(data))
	if kpb == nil {
		return nil, errors.New("Private key did not contain PEM formatted block")
	}

	key, err := x509.ParsePKCS1PrivateKey(kpb.Bytes)
	if err != nil {
		return nil, errors.New(err.Error() + "Parsing private key")
	}

	return key, nil
}
