package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

func generateCertificate(ca *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate,
	*rsa.PrivateKey, *bytes.Buffer, *bytes.Buffer, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Alex Socks5"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(30 * 24 * time.Hour),

		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	if ca == nil {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		template.BasicConstraintsValid = true
		ca = &template
		caKey = priv
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, ca, &priv.PublicKey, caKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	var cert bytes.Buffer
	if err := pem.Encode(&cert, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return nil, nil, nil, nil, err
	}
	var key bytes.Buffer
	if err := pem.Encode(&key, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, nil, nil, err
	}
	return &template, priv, &cert, &key, nil
}

func GenerateCertificates() (tls.Certificate, *x509.CertPool, tls.Certificate, error) {
	ca, caKey, caCert, _, err := generateCertificate(nil, nil)
	if err != nil {
		return tls.Certificate{}, nil, tls.Certificate{}, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert.Bytes()) {
		return tls.Certificate{}, nil, tls.Certificate{}, errors.New("failed to add ca cert to pool")
	}
	_, _, cert, key, err := generateCertificate(ca, caKey)
	if err != nil {
		return tls.Certificate{}, nil, tls.Certificate{}, err
	}
	serverCert, err := tls.X509KeyPair(cert.Bytes(), key.Bytes())
	if err != nil {
		return tls.Certificate{}, nil, tls.Certificate{}, err
	}
	_, _, cert, key, err = generateCertificate(ca, caKey)
	if err != nil {
		return tls.Certificate{}, nil, tls.Certificate{}, err
	}
	clientCert, err := tls.X509KeyPair(cert.Bytes(), key.Bytes())
	if err != nil {
		return tls.Certificate{}, nil, tls.Certificate{}, err
	}
	return serverCert, pool, clientCert, nil
}
