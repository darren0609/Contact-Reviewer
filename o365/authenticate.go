package o365

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// getCreds - will read private credentials from text file and return them for use later within the routers.
func getCreds(filepath string) (string, string, error) {
	var err error
	var id, secret string
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		id = scanner.Text()
	}
	if scanner.Scan() {
		secret = scanner.Text()
	}
	if id[0] == '*' || secret[0] == '*' {
		err := errors.New("Missing Configuration: _PRIVATE.txt needs to be edited to add client ID and secret")
		return "", "", err
	}
	return id, secret, err
}

// Handler for login route
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// guid is created using the Google UUID library - it is set randomly
	// guid used to only accept initiated logins, compared after response later

	if guid == "" {
		tmpguid, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		} else {
			guid = tmpguid.String()
		}
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"User.Read", "Contacts.Read", "Mail.Send"},
		RedirectURL:  "http://localhost:8080/login",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
	}
	var code string
	code = r.URL.Query().Get("code")
	if len(code) == 0 {
		// Redirect user to consent page to ask for permission
		// for the scopes specified above.
		authurl := conf.AuthCodeURL(guid, oauth2.AccessTypeOffline)
		http.Redirect(w, r, authurl, http.StatusSeeOther)
		// fmt.Printf(authurl)
		return
	}
	// Before calling Exchange, be sure to validate FormValue("state").
	if r.FormValue("state") != guid {
		log.Fatal("State has been messed with, end authentication")
		// reset state to prevent re-use
		guid = ""
	}
	ctx := context.Background()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	client = conf.Client(ctx, token)

	// Grab credentials so we can use them in displaying form detail going forward
	res, err := client.Get("https://graph.microsoft.com/v1.0/me")

	if err != nil {
		log.Println("Failed to get user/me:", err)
		return
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		log.Println("Failed to parse user data:", err)
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return
}
