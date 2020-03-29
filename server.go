package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type Message struct {
	Address    string `json:"address"`
	UserAgent  string `json:"user_agent"`
	DomainName string `json:"domain_name"`
}

func main() {
	hostPtr := flag.String("host", "0.0.0.0", "Host address to listen on")
	portPtr := flag.Int("port", 8080, "Port to listen on")
	httpsPtr := flag.Bool("https", false, "Serve over HTTPS")
	pubPtr := flag.String("public", "", "Public key (.pem) - Requires https flag")
	keyPtr := flag.String("private", "", "Private key (.pem) - Requires https flag")
	flag.Parse()

	socket := fmt.Sprintf("%s:%d", *hostPtr, *portPtr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	server := http.Server{
		Addr:    socket,
		Handler: mux,
	}

	log.Printf("Listening on %s", socket)

	if *httpsPtr {
		if *pubPtr == "" || *keyPtr == "" {
			fmt.Fprintf(os.Stderr, "Error: HTTPS flag provided without valid public and private options")
			os.Exit(1)
		}

		log.Fatal(server.ListenAndServeTLS(*pubPtr, *keyPtr))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	var addr, err = extractAddressFromSocket(request.RemoteAddr)

	if err != nil {
		fmt.Fprintf(writer, "Ip address was whack.")
		return
	}
	realIP := request.Header.Get("X-Real-IP")
	if realIP != "" {
		addr = realIP
	}

	var name = lookup(addr)

	data := Message{
		Address:    addr,
		UserAgent:  request.UserAgent(),
		DomainName: name,
	}

	log.Printf("%s -- %s -- %s", data.Address, data.DomainName, data.UserAgent)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

func lookup(addr string) string {
	names, err := net.LookupAddr(addr)

	if err != nil {
		return "Error resolving address"
	}

	var result string
	if len(names) == 0 {
		result = "No names found"
	} else if len(names) == 1 {
		result = names[0]
	} else {
		result = fmt.Sprintf("%s (+%d other names)", names[0], len(names))
	}

	return result
}

func extractAddressFromSocket(addr string) (string, error) {
	ip, _, err := net.SplitHostPort(addr)
	if err == nil {
		return ip, nil
	}

	ip2 := net.ParseIP(addr)
	if ip2 == nil {
		return "", errors.New("invalid IP")
	}

	return ip2.String(), nil
}
