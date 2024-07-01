package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"

	"github.com/things-go/go-socks5"
)

func ListenAndServeTLS(server *socks5.Server, addr string, config *tls.Config) error {
	l, err := tls.Listen("tcp", addr, config)
	if err != nil {
		return err
	}
	return server.Serve(l)
}

func main() {
	port := flag.String("port", "18881", "listening port")
	certFile := flag.String("cert", "example.crt", "certificate PEM file")
	keyFile := flag.String("key", "example.key", "key PEM file")
	user := flag.String("user", "", "user name")
	pass := flag.String("pass", "", "password")
	flag.Parse()

	if *user == "" || *pass == "" {
		log.Printf("You need to provide user and pass.")
		return
	}

	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithCredential(socks5.StaticCredentials{*user: *pass}),
	)
	addr := ":" + *port
	log.Printf("Listening on %s with TLS.\n", addr)
	if err := ListenAndServeTLS(server, addr, config); err != nil {
		panic(err)
	}
}
