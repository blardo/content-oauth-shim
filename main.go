package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	flag "github.com/ogier/pflag"
	"github.com/skuid/content-oauth-shim/config"
	"github.com/skuid/content-oauth-shim/middlewares"
	"github.com/skuid/content-oauth-shim/views"
	"strings"
)

const versionString = "0.0.1"

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, part := range parts {
		*i = append(*i, part)
	}
	return nil
}

var emailDomains arrayFlags

var version = flag.Bool("version", false, "print version and exit")

var port = flag.Int("port", 3000, "The port to listen on")

// TODO Make session store pluggable
var redisConnection = flag.String("redis", "192.168.99.100:6379", "The redis connection param to use")
var sessionSecret = flag.String("secret", "", "The session secret to use")
var assetRoot = flag.String("assets", "./assets", "The asset root to use")

var redirectHost = flag.String("redirect-host", "", "The scheme://host[:port] for Google to return to")

var tlsKey = flag.String("tls-key", "", "The path to the TLS key to use. tls-key and tls-cert must be used to serve TLS")
var tlsCert = flag.String("tls-cert", "", "The path to the TLS cert to use")

var clientID = flag.String("client-id", os.Getenv("CLIENT_ID"), "The Client ID to use")
var clientSecret = flag.String("client-secret", os.Getenv("CLIENT_SECRET"), "The Client Secret to use")

var appName = flag.String("app-name", "Oauth Application", "The application name to use")

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	domainHelp := "Email domain that authenticated users must match. Can be " +
		"specified multiple times, or once separating domains with a comma. \n" +
		"        Not listing domains will result in any google user being allowed " +
		"to authenticate."
	flag.Var(&emailDomains, "domain", domainHelp)
	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", os.Args[0], versionString)
		os.Exit(0)
	}

	sc, err := config.NewServerConfig(
		*sessionSecret,
		config.RedisConnection(*redisConnection),
		config.Port(*port),
	)
	if err != nil {
		panic(err)
	}

	r, err := views.NewContentRouter(
		*assetRoot,
		*clientID,
		*clientSecret,
		views.SessionStore(sc.RedisStore),
		views.AuthorizedDomains(emailDomains),
		views.RedirectHost(*redirectHost),
		views.ApplicationName(*appName),
	)
	if err != nil {
		panic(err)
	}

	r.Router.HandleFunc("/oauth/callback", r.OauthCallback).Methods("POST")
	r.Router.HandleFunc("/login", r.LoginView).Methods("GET")
	r.Router.HandleFunc("/logout", r.LogoutView).Methods("GET")
	//r.Router.PathPrefix("/").HandlerFunc(r.Assets)
	r.Router.PathPrefix("/").HandlerFunc(r.Content)
	http.Handle("/", r.Router)
	log.Printf("Listening on %s", sc.Hostport())

	if len(*tlsKey) > 0 && len(*tlsCert) > 0 {
		log.Fatal(
			http.ListenAndServeTLS(sc.Hostport(), *tlsCert, *tlsKey,
				middlewares.Apply(r.Router, middlewares.Logging()),
			),
		)
	} else {
		log.Fatal(
			http.ListenAndServe(
				sc.Hostport(),
				middlewares.Apply(r.Router, middlewares.Logging()),
			),
		)
	}
}
