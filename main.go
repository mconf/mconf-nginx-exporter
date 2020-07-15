package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var ips []string

func exporterHandler(w http.ResponseWriter, r *http.Request) {
	resp := ""

	for _, ip := range ips {
		u := "http://" + ip + "/metrics"
		r, err := http.Get(u)
		if err != nil {
			fmt.Println(err)
		} else {
			body := r.Body
			bytes, err := ioutil.ReadAll(body)
			if err != nil {
				fmt.Println(err)
			} else {
				resp += string(bytes)
			}
		}
	}

	w.Write([]byte(resp))
}

func main() {
	ipsStr := os.Getenv("NGINX_INTERNAL_IPS")
	if ipsStr == "" {
		fmt.Println("No ips to run with, exiting")
		return
	}
	ips = strings.Split(ipsStr, ",")
	fmt.Println("Running with ips:", ips)

	http.HandleFunc("/metrics", exporterHandler)

	port := os.Getenv("NGINX_EXPORTER_PORT")
	if port == "" {
		port = "8888"
	}
	fmt.Printf("Listening at port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

/*
	hmg: 12GB used
*/
