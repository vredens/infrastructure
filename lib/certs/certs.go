package certs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Locations []string
}

type Certs struct {
	certs []string
	root  *x509.CertPool
}

func New(config Config) Certs {
	certs := Certs{}
	certs.loadCerts(config.Locations)

	return certs
}

func (certs Certs) RootCAs() *x509.CertPool {
	return certs.root
}

func (certs Certs) NewTLSClientConfig() *tls.Config {
	return &tls.Config{
		RootCAs: certs.RootCAs(),
	}
}

func (certs *Certs) loadCerts(locations []string) {
	for i := range locations {
		ok, err := certs.loadFromLocation(locations[i])
		if err != nil {
			continue
		}
		if ok {
			return
		}
	}
}

func (certs *Certs) loadFromLocation(path string) (ok bool, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to check folder; %w", err)
	}
	if !info.IsDir() {
		return false, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("failed to read folder; %w", err)
	}
	if len(files) == 0 {
		return false, nil
	}

	if certs.root, err = x509.SystemCertPool(); err != nil {
		certs.root = x509.NewCertPool()
	}

	for i := range files {
		info, err = files[i].Info()
		if err != nil {
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}
		if info.Name() == "ca.pem" || strings.HasSuffix(info.Name(), ".ca.pem") {
			data, err := os.ReadFile(filepath.Join(path, info.Name()))
			if err != nil {
				continue
			}
			if len(data) == 0 {
				continue
			}
			if certs.root.AppendCertsFromPEM(data) {
				ok = true
				certs.certs = append(certs.certs, info.Name())
			}
			continue
		}
	}

	if !ok {
		certs.root = nil
		certs.certs = nil
	}
	return ok, nil
}
