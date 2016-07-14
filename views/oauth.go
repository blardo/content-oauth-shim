package views

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type userInfo struct {
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
}

type responseError struct {
	Domain       string `json:"domain"`
	Reason       string `json:"reason"`
	Message      string `json:"message"`
	LocationType string `json:"locationType"`
	Location     string `json:"location"`
}

type responseErrorInfo struct {
	Errors  []*responseError `json:"errors"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
}

type emailResponse struct {
	Data   *userInfo `json:"data"`
	Errors *responseErrorInfo
}

var scopes = []string{"profile", "email"}

type googleResponse struct {
	Iss          string `json:"iss"`
	Sub          string `json:"sub"`
	Azp          string `json:"azp"`
	Aud          string `json:"aud"`
	Iat          string `json:"iat"`
	Exp          string `json:"exp"`
	HostedDomain string `json:"hd"`

	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified,string"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

func (router *ContentRouter) verifyToken(token string) bool {
	uri := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token

	response, err := http.Get(uri)
	defer response.Body.Close()
	if err != nil {
		log.Printf("Error validating token %s", err.Error())
		return false
	}
	if response.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(response.Body)
		log.Printf("Error in response %s", string(data))
		return false
	}
	data := &googleResponse{}
	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		log.Printf("Error decoding json %s", err.Error())
		return false
	}

	if len(data.Email) == 0 {
		log.Printf("Email not provided in response %#v", data)
		return false
	}
	_, ok := router.AuthorizedDomains[data.HostedDomain]
	if ok && data.EmailVerified {
		log.Printf("Successfully authenticated %s: <%s>", data.Name, data.Email)
		return true
	}
	log.Printf("Email domain doesn't match %s", data.Email)
	return false
}

// OauthCallback - The callback view for Google to hit
func (router *ContentRouter) OauthCallback(w http.ResponseWriter, r *http.Request) {
	accessToken := r.FormValue("idtoken")
	authorized := true
	if len(router.AuthorizedDomains) > 0 {
		authorized = router.verifyToken(accessToken)
	}

	if authorized {
		session, err := router.SessionStore.Get(r, "content-sid")
		if err != nil {
			log.Printf("Error getting session %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "401 - Not Authorized", http.StatusUnauthorized)
}

var loginTemplate = template.Must(template.New("login").Parse(login_template))

// LoginView - The Login View
func (router *ContentRouter) LoginView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	session, err := router.SessionStore.Get(r, router.SessionName)
	if err != nil {
		log.Printf("Error getting session %s", err.Error())
	}
	if !session.IsNew {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	router.CallbackHost.Path = "/oauth/callback"

	// Fill in the missing fields in index.html
	context := struct {
		ApplicationName string
		ClientID        string
		CallbackURL     string
		Scopes          []string
	}{
		router.ApplicationName,
		router.ClientID,
		router.CallbackHost.String(),
		scopes,
	}

	// Render and serve the HTML
	err = loginTemplate.Execute(w, context)
	if err != nil {
		log.Println("error rendering template:", err)
		http.Error(w, "Error Rendering Template", http.StatusInternalServerError)
		return
	}
}

// LogoutView removes the user's session
func (router *ContentRouter) LogoutView(w http.ResponseWriter, r *http.Request) {
	session, err := router.SessionStore.Get(r, "content-sid")
	if err != nil {
		log.Printf("Error getting session: %s", err.Error())
		http.Error(w, "Error getting session", http.StatusInternalServerError)
		return
	}
	if !session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Logged Out`))
}
