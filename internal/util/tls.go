// Copyright 2023 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

// CreateCACertPool creates a CA cert pool gives a CA cert bundle
func CreateCACertPool(caCertBundlePath string) (*x509.CertPool, error) {
	caCertBundle, err := os.ReadFile(caCertBundlePath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBundle)
	return caCertPool, nil
}

// IsTLSEnabledFn returns true if TLS is enabled and false otherwise.
type IsTLSEnabledFn func() bool

// KeyPair is a collection of paths one for the certificate and another for the key.
// This is used to configure certificate-key pair when configuring TLS config.
type KeyPair struct {
	// CertPath is the path to the certificate
	CertPath string
	// KeyPath is the path to the private key
	KeyPath string
}

// CreateTLSConfig creates a TLS Config to be used for TLS communication.
func CreateTLSConfig(tlsEnabledFn IsTLSEnabledFn, serverName, caCertPath string, keyPair *KeyPair) (*tls.Config, error) {
	tlsConf := tls.Config{}
	if !tlsEnabledFn() {
		tlsConf.InsecureSkipVerify = true
		return &tlsConf, nil
	}

	caCertPool, err := CreateCACertPool(caCertPath)
	if err != nil {
		return nil, err
	}
	tlsConf.RootCAs = caCertPool
	tlsConf.ServerName = serverName
	if keyPair != nil {
		certificate, err := tls.LoadX509KeyPair(keyPair.CertPath, keyPair.KeyPath)
		if err != nil {
			return nil, err
		}
		tlsConf.Certificates = []tls.Certificate{certificate}
	}
	return &tlsConf, nil
}
