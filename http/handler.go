package http

import (
	"crypto"
	"fmt"
	"github.com/davepgreene/slackmac/store"
	"github.com/davepgreene/slackmac/utils"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
	"net/http"
)

//var statsMiddleware *stats.Stats

// Handler returns an http.Handler for the API.
func Handler() error {
	storeConf := viper.GetStringMapString("store")
	dataStore, err := store.CreateStore(storeConf)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	stats := stats.New()

	// Add middleware handlers
	n := negroni.New()

	// Add recovery handler that logs
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n.Use(recovery)

	if viper.GetBool("log.requests") {
		n.Use(negronilogrus.NewCustomMiddleware(utils.GetLogLevel(), utils.GetLogFormatter(), "requests"))
	}

	n.Use(stats)

	// Collect some metrics about incoming and active requests
	n.Use(negroni.HandlerFunc(metricsMiddleware))

	// Correlation id middleware
	if viper.GetBool("correlation.enable") {
		n.Use(negroni.HandlerFunc(correlationMiddleware(viper.GetStringMap("correlation"))))
	}

	// Validate headers and timestamp
	skew := viper.GetDuration("slack.skew")
	n.Use(negroni.HandlerFunc(timestamp(skew)))

	// Validate the signature
	algorithm := viper.Get("slack.algorithm")
	n.Use(negroni.HandlerFunc(signature(dataStore, algorithm.(crypto.Hash))))

	// All checks passed, forward the request
	forwardConn := fmt.Sprintf("%s%s:%d", viper.GetString("service.protocol"), viper.GetString("service.hostname"), viper.GetInt("service.port"))
	p := proxy(forwardConn)
	r.PathPrefix("/").Handler(p)

	n.UseHandler(r)

	// Set up connection
	conn := fmt.Sprintf("%s:%d", viper.GetString("listen.bind"), viper.GetInt("listen.port"))
	log.Infof("Listening on %s", conn)

	// Bombs away!
	go func() {
		m := mux.NewRouter()
		statsMiddleware := negroni.New()

		m.HandleFunc("/stats", newAdminHandler(stats).ServeHTTP)

		if viper.GetBool("admin.log") {
			statsMiddleware.Use(negronilogrus.NewCustomMiddleware(utils.GetLogLevel(), utils.GetLogFormatter(), "admin.requests"))
		}

		statsMiddleware.UseHandler(m)

		adminConn := fmt.Sprintf("%s:%d", viper.GetString("admin.bind"), viper.GetInt("admin.port"))
		log.Fatal(server(adminConn, statsMiddleware).ListenAndServe())
	}()
	return server(conn, n).ListenAndServe()
}

func server(conn string, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    conn,
		Handler: handler,
	}

	srv.SetKeepAlivesEnabled(false)

	return srv
}
