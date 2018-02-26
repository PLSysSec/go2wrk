package https

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

// SetTLS sets the TLS for a server. We want to look into this more.
func SetTLS(disableKeepAlives, insecure bool, certFile, keyFile, caFile string) *http.Transport {

	var tlsConfig *tls.Config

	if insecure {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig, DisableKeepAlives: disableKeepAlives}

	return transport
}
