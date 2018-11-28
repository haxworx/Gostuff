package main;

import "log"
import "context"
import "io/ioutil"
import "net/http"
import "net/url"
import "golang.org/x/oauth2"
import "golang.org/x/oauth2/facebook"
import "golang.org/x/oauth2/google"


var (
	oauthFacebookConf = &oauth2.Config {
		ClientID:	"CLIENT_ID",
		ClientSecret:	"CLIENT_SECRET",
		RedirectURL:	"https://localhost:8080/facebook/callback",
		Scopes:		[]string{"public_profile"},
		Endpoint:	facebook.Endpoint,
	}

	oauthGoogleConf = &oauth2.Config {
                ClientID:       "CLIENT_ID",
                ClientSecret:   "CLIENT_SECRET",
                RedirectURL:    "http://localhost:8080/google/callback",
                Scopes:         []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
                Endpoint:       google.Endpoint,

	}

	oauthStateString = "randomstring"
)

func authGoogle(w http.ResponseWriter, r *http.Request) {
	url := oauthGoogleConf.AuthCodeURL(oauthStateString)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state"); if state != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := oauthGoogleConf.Exchange(context.Background(), code); if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" +
			      url.QueryEscape(token.AccessToken)); if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body); if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	log.Printf("Google %s\n", string(response));

	http.Redirect(w, r, "/", http.StatusOK)
}

func authFacebook(w http.ResponseWriter, r *http.Request) {

	url := oauthFacebookConf.AuthCodeURL(oauthStateString)


	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func authFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state"); if state != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := oauthFacebookConf.Exchange(oauth2.NoContext, code); if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
			     url.QueryEscape(token.AccessToken))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect);
		return
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body); if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	log.Printf("Hello %s\n", string(response))

	http.Redirect(w, r, "/", http.StatusOK)
}

func main() {

	http.HandleFunc("/google", authGoogle)
	http.HandleFunc("/google/callback", authGoogleCallback)

	http.HandleFunc("/facebook", authFacebook)
	http.HandleFunc("/facebook/callback", authFacebookCallback)

	err := http.ListenAndServe(":8080", nil); if err != nil {
		log.Fatal(err)
	}

}
