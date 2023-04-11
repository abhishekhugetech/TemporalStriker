package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"go.uber.org/zap"
)

func GetTLSConfig(hostPort string, logger *zap.Logger) (*tls.Config, error) {
	host, _, parseErr := net.SplitHostPort(hostPort)
	if parseErr != nil {
		return nil, fmt.Errorf("unable to parse hostport properly: %+v", parseErr)
	}

	caCertData := GetEnvOrDefaultString(logger, "TLS_CA_CERT_DATA", "")
	clientCertData := GetEnvOrDefaultString(logger, "TLS_CLIENT_CERT_DATA", "")
	clientCertPrivateKeyData := GetEnvOrDefaultString(logger, "TLS_CLIENT_CERT_PRIVATE_KEY_DATA", "")
	caCertFile := GetEnvOrDefaultString(logger, "TLS_CA_CERT_FILE", "")
	clientCertFile := GetEnvOrDefaultString(logger, "TLS_CLIENT_CERT_FILE", "")
	clientCertPrivateKeyFile := GetEnvOrDefaultString(logger, "TLS_CLIENT_CERT_PRIVATE_KEY_FILE", "")
	enableHostVerification := GetEnvOrDefaultBool(logger, "TLS_ENABLE_HOST_VERIFICATION", false)

	caBytes, err := getTLSBytes(caCertFile, caCertData)
	if err != nil {
		return nil, err
	}

	certBytes, err := getTLSBytes(clientCertFile, clientCertData)
	if err != nil {
		return nil, err
	}

	keyBytes, err := getTLSBytes(clientCertPrivateKeyFile, clientCertPrivateKeyData)
	if err != nil {
		return nil, err
	}

	var cert *tls.Certificate
	var caPool *x509.CertPool

	if len(certBytes) > 0 {
		clientCert, err := tls.X509KeyPair(certBytes, keyBytes)
		if err != nil {
			return nil, err
		}
		cert = &clientCert
	}

	if len(caBytes) > 0 {
		caPool = x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caBytes) {
			return nil, errors.New("unknown failure constructing cert pool for ca")
		}
	}

	// If we are given arguments to verify either server or client, configure TLS
	if caPool != nil || cert != nil {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: !enableHostVerification,
			ServerName:         host,
		}
		if caPool != nil {
			tlsConfig.RootCAs = caPool
		}
		if cert != nil {
			tlsConfig.Certificates = []tls.Certificate{*cert}
		}

		return tlsConfig, nil
	}

	return nil, nil

}

func getTLSBytes(certFile string, certData string) ([]byte, error) {
	var bytes []byte
	var err error

	if certFile != "" && certData != "" {
		return nil, errors.New("cannot specify both file and Base-64 encoded version of same field")
	}

	if certFile != "" {
		bytes, err = ioutil.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
	} else if certData != "" {
		bytes, err = base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return nil, err
		}
	}

	return bytes, err
}
