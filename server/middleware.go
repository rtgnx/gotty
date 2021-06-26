package server

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
)

func (server *Server) wrapLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &logResponseWriter{w, 200}
		handler.ServeHTTP(rw, r)
		log.Printf("%s %d %s %s", r.RemoteAddr, rw.status, r.Method, r.URL.Path)
	})
}

func (server *Server) wrapHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo add version
		w.Header().Set("Server", "GoTTY")
		handler.ServeHTTP(w, r)
	})
}

func (server *Server) wrapBasicAuth(handler http.Handler, credential string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(token) != 2 || strings.ToLower(token[0]) != "basic" {
			w.Header().Set("WWW-Authenticate", `Basic realm="GoTTY"`)
			http.Error(w, "Bad Request", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(token[1])
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if credential != string(payload) {
			w.Header().Set("WWW-Authenticate", `Basic realm="GoTTY"`)
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		log.Printf("Basic Authentication Succeeded: %s", r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}

func (server *Server) wrapForwardBasicAuth(handler http.Handler, forwardAuthServer string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(token) != 2 || strings.ToLower(token[0]) != "basic" {
			w.Header().Set("WWW-Authenticate", `Basic realm="GoTTY"`)
			http.Error(w, "Bad Request", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(token[1])
		if err != nil || !strings.Contains(string(payload), ":") {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("Authenticatig with %s", payload)
		credentials := strings.Split(string(payload), ":")

		if len(credentials) != 2 {
			http.Error(w, "Invalid Credentials Supplied", http.StatusBadRequest)
		}

		client := &http.Client{}
		req, _ := http.NewRequest("GET", forwardAuthServer, nil)

		req.SetBasicAuth(credentials[0], credentials[1])

		res, err := client.Do(req)

		if err != nil {

			log.Println(err.Error())

			w.Header().Set("WWW-Authenticate", `Basic realm="GoTTY"`)
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		if res.StatusCode == 200 {
			for k, v := range res.Header {
				log.Printf("Updating Headers %s = %v\n", k, v)
				r.Header.Set(k, strings.Join(v, ","))
				// Replace - with _ for Environment variables to be accessible
				k = strings.Replace(k, "-", "_", -1)
				os.Setenv(k, strings.Join(v, ","))
			}
			log.Printf("Basic Authentication Succeeded: %s", r.RemoteAddr)
			handler.ServeHTTP(w, r)
		}

		if res.StatusCode != 200 {
			log.Printf("Forward Auth Response: %d\n", res.StatusCode)

			w.Header().Set("WWW-Authenticate", `Basic realm="GoTTY"`)
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
	})
}
