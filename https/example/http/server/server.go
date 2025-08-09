package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1) Nạp CA để xác thực client-cert.pem
	caCertPEM, err := os.ReadFile("../../../cert/ca-cert.pem")
	if err != nil {
		log.Fatal("read ca-cert.pem:", err)
	}
	clientCAPool := x509.NewCertPool()
	if !clientCAPool.AppendCertsFromPEM(caCertPEM) {
		log.Fatal("append CA failed")
	}

	// 2) Nạp chứng chỉ + key của server
	serverCert, err := tls.LoadX509KeyPair("../../../cert/server-cert.pem", "../../../cert/server-key.pem")
	if err != nil {
		log.Fatal("load server cert/key:", err)
	}

	// 3) Cấu hình TLS: bắt buộc client đưa cert và verify với CA
	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    clientCAPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // mTLS
		MinVersion:   tls.VersionTLS12,
		// ServerName trong client phải khớp SAN của server-cert.
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// Lấy thông tin client cert (đã verify)
		state := r.TLS
		var who string
		if state != nil && len(state.PeerCertificates) > 0 {
			who = state.PeerCertificates[0].Subject.CommonName
		} else {
			who = "unknown-client"
		}
		io.WriteString(w, fmt.Sprintf("Hello, %s! You’re mTLS-authenticated.\n", who))
	})

	srv := &http.Server{
		Addr:      ":8443",
		Handler:   mux,
		TLSConfig: tlsCfg,
	}

	log.Println("Server listening on https://localhost:8443")
	// ListenAndServeTLS sẽ dùng Certificates trong tlsCfg, nên để ""/""
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
