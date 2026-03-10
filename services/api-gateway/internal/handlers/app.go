package handlers

import "net/http"

func RootHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("API-Gateway v1.0.0"))
}

func HealthHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK"))
}
