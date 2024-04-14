package utilities

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func SetupHttpsServer(router http.Handler, cert, key, httpsPort, httpPort string, requestTimeout time.Duration) (server *http.Server) {
	if FileExists(cert) && FileExists(key) {
		server = &http.Server{
			Addr:         httpsPort,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		go func() {
			log.Println("Server starts at port ", httpsPort)
			if err := server.ListenAndServeTLS(cert, key); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()
	} else {
		server = &http.Server{
			Addr:         httpPort,
			Handler:      router,
			ReadTimeout:  requestTimeout * time.Second,
			WriteTimeout: requestTimeout * time.Second,
		}
		go func() {
			log.Println("Server starts at port ", httpPort)
			if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()
	}

	return server
}
