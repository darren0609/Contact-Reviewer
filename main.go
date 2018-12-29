package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var (
	clientID     string
	clientSecret string
	guid         string
	user         User
	client       *http.Client
)

// D is Data for template
type D map[string]interface{}

// User represents a Microsoft Graph user
type User struct {
	Username  string `json:"displayName"`
	Email     string `json:"mail"`
	Surname   string `json:"surname"`
	Givenname string `json:"givenname"`
}

// ContactData is the contact detail
type ContactData struct {
	//Etag                 string   `json:"odata.etag"`
	//ID                   string   `json:"id"`
	//Createddatetime      string   `json:"createdDateTime"`
	//Lastmodifieddatetime string   `json:"lastModifiedDateTime"`
	//Changekey            string   `json:"changeKey"`
	Displayname     string         `json:"displayName"`
	Officephone     []string       `json:"businessPhones"`
	Homephone       []string       `json:"homePhones"`
	Emailaddresses  []EmailAddress `json:"emailAddresses"`
	Givenname       string         `json:"givenName"`
	Mobilephone     string         `json:"mobilePhone"`
	Surname         string         `json:"surname"`
	Businessaddress []Address      `json:"businessAddress"`
}

// Address is the sub-detail within addresses
type Address struct {
	Street          string `json:"street"`
	City            string `json:"city"`
	State           string `json:"state"`
	Countryorregion string `json:"countryOrRegion"`
	Postalcode      string `json:"postalCode"`
}

// ContactHeader is the contact header and detail
type ContactHeader struct {
	//Context  string        `json:"odata.context"`
	//NextLink string        `json:"odata.nextLink"`
	Contacts []ContactData `json:"value"`
}

// Body is generally the body of the returned HTML when executing.
type Body struct {
	ContentType string
	Content     string
}

// EmailAddress is the Email Address of the user
type EmailAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Recipient is the Receipient of the email you wish to send
type Recipient struct {
	EmailAddress EmailAddress
}

// Message is the specifics of the message being sent
type Message struct {
	Subject      string
	Body         Body
	ToRecipients []Recipient
}

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
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	client = conf.Client(ctx, tok)

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

func main() {
	var err error
	// Configure API ClientID/Secret from configuration file
	configFile := "init/private.txt"
	clientID, clientSecret, err = getCreds(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", loginHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/home", homeHandler)
	r.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", r)

	//fmt.Println("Client ID : ", clientID)
	//fmt.Println("Client Secret : ", clientSecret)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var endpointURL string
	var results ContactHeader
	var err error

	// Use OData query parameters to control the results
	// - Only first 10 results returned
	// - Only return the GivenName, Surname, and EmailAddresses fields
	// - Sort the results by the GivenName field in ascending order
	//query_parameters := "$top=50",
	//	"$select" : "givenName,surname,emailAddresses",
	//	"$orderby": "givenName ASC",
	//}

	// Post the message to the graph API endpoint for sending email

	//$orderby=givenName ASC$top=50"
	//endpointURL := "https://graph.microsoft.com/v1.0/me/contacts"

	endpointURL = "https://graph.microsoft.com/v1.0/me/contacts"

	if r.FormValue("search") != "" {
		endpointURL = endpointURL + "?$search=" + r.FormValue("search")
	}

	fmt.Println("Before I check the formvalue: ", r.FormValue("sortBy"))

	if r.FormValue("sortBy") != "" {

		endpointURL = endpointURL + "?$orderby=" + r.FormValue("sortBy")
		fmt.Println("Made it inside the sortBy setting")
	}

	fmt.Println("EndpointURL: ", endpointURL)
	res, err := client.Get(endpointURL)
	if err != nil {
		fmt.Println("Error in get contacts:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &results)
	if err != nil {
		log.Println("Failed to UNMARSHAL user data:", err)
	}

	fmt.Println("Results are : ", results)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler for sendmail route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	var fullcont ContactHeader
	var endpointURL string
	var err error

	// Use OData query parameters to control the results
	// - Only first 10 results returned
	// - Only return the GivenName, Surname, and EmailAddresses fields
	// - Sort the results by the GivenName field in ascending order
	//query_parameters := "$top=50",
	//	"$select" : "givenName,surname,emailAddresses",
	//	"$orderby": "givenName ASC",
	//}

	// Post the message to the graph API endpoint for sending email

	//$orderby=givenName ASC$top=50"
	//endpointURL := "https://graph.microsoft.com/v1.0/me/contacts"

	endpointURL = "https://graph.microsoft.com/v1.0/me/contacts"

	if r.FormValue("search") != "" {
		endpointURL = endpointURL + "?$search=" + r.FormValue("search")
	}

	res, err := client.Get(endpointURL)
	if err != nil {
		fmt.Println("Error in get contacts:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &fullcont)
	if err != nil {
		log.Println("Failed to UNMARSHAL user data:", err)
	}

	// Parse template for response to app client
	t2, err := template.ParseFiles("tpl/contacts.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t2.Execute(w, D{
		"me":          user,
		"contact":     fullcont,
		"showSuccess": false,
		"showError":   false,
	})

	if err != nil {
		fmt.Println("Error executing template pass with data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
