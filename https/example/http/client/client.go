package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1) Nạp CA để verify server-cert.pem
	caCertPEM, err := os.ReadFile("../../../cert/ca-cert.pem")
	if err != nil {
		log.Fatal("read ca-cert.pem:", err)
	}
	rootPool := x509.NewCertPool()
	if !rootPool.AppendCertsFromPEM(caCertPEM) {
		log.Fatal("append CA failed")
	}

	// 2) Nạp client cert + key để server verify
	clientCert, err := tls.LoadX509KeyPair("../../../cert/client-cert.pem", "../../../cert/client-key.pem")
	if err != nil {
		log.Fatal("load client cert/key:", err)
	}

	tlsCfg := &tls.Config{
		RootCAs:      rootPool,                      // verify server
		Certificates: []tls.Certificate{clientCert}, // present client cert
		MinVersion:   tls.VersionTLS12,
		ServerName:   "localhost",
	}

	tr := &http.Transport{TLSClientConfig: tlsCfg}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatal("request failed:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Status: %s\nBody:\n%s", resp.Status, string(body))
}
