package https

import (
	"crypto/x509"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"log"
)

func SetTLS(disable_keep_alives, insecure bool, cert_file, key_file, ca_file string) *http.Transport {

	var tls_config *tls.Config

	if insecure {
		tls_config = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(cert_file, key_file)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		ca_cert, err := ioutil.ReadFile(ca_file)
		if err != nil {
			log.Fatal(err)
		}
		ca_cert_pool := x509.NewCertPool()
		ca_cert_pool.AppendCertsFromPEM(ca_cert)

		// Setup HTTPS client
		tls_config = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      ca_cert_pool,
		}
		tls_config.BuildNameToCertificate()
	}

	transport := &http.Transport{TLSClientConfig: tls_config, DisableKeepAlives: disable_keep_alives}

	return transport
}