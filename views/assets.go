package views

import (
	"log"
	"net/http"
	"path"
)

// Content serves the files out of the router.AssetRoot
func (router *ContentRouter) Content(w http.ResponseWriter, r *http.Request) {
	session, err := router.SessionStore.Get(r, router.SessionName)
	if err != nil {
		log.Printf("Error creating session: %s", err.Error())
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	if session.IsNew {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	file := path.Join(router.AssetRoot, r.URL.Path[1:])
	http.ServeFile(w, r, file)
}
