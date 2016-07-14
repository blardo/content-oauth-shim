package views

import (
	"log"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// ContentRouter is a structure for passing information to views
type ContentRouter struct {
	Router            *mux.Router     // The mux router
	AssetRoot         string          // The asset root to serve
	SessionStore      sessions.Store  // The session store to use
	AuthorizedDomains map[string]bool // The domains authorized to view content
	CallbackHost      *url.URL        // The scheme://host[:port] for Google to return to
	ClientID          string          // The Oauth Client Id
	ClientSecret      string          // the Oauth Client Secret
	ApplicationName   string          // Name of the Application
	SessionName       string          // Name of the session string
}

// RouterOption is a type for setting options on the ContentRouter
type RouterOption func(*ContentRouter)

// SessionStore adds a session store to the router
func SessionStore(store sessions.Store) RouterOption {
	return func(ro *ContentRouter) {
		ro.SessionStore = store
	}
}

// AuthorizedDomains adds authorized domains to the router
func AuthorizedDomains(domains []string) RouterOption {
	return func(ro *ContentRouter) {
		for _, d := range domains {
			ro.AuthorizedDomains[d] = true
		}
	}
}

// RedirectHost sets a the domain for the oauth callback URL
func RedirectHost(host string) RouterOption {
	return func(ro *ContentRouter) {
		uri, err := url.Parse(host)
		if err != nil {
			log.Fatal(err)
		}
		ro.CallbackHost = uri
	}
}

// ApplicationName sets a the application name
func ApplicationName(name string) RouterOption {
	return func(ro *ContentRouter) {
		ro.ApplicationName = name
	}
}

// NewContentRouter returns a router with the specified option set
func NewContentRouter(assetRoot, clientID, clientSecret string, options ...RouterOption) (*ContentRouter, error) {
	dr := &ContentRouter{
		Router:            mux.NewRouter(),
		AssetRoot:         assetRoot,
		AuthorizedDomains: map[string]bool{},
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		SessionName:       "content-sid",
	}

	for i := range options {
		if options[i] != nil {
			options[i](dr)
		}
	}
	return dr, nil
}
