package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

var (
	redirectURI  = "https://hankdoupe.com/ttrack.html"
	clientID     = ""
	clientSecret = ""
	// AuthURL is where the user should be directed to login to FreshBooks.
	AuthURL = ""
)

// Credentials contains the data from a successful authentication
// flow.
type Credentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

// Access contains the data for retrieving Credentials.
type Access struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

// Refresh contains the data for refreshing expired Credentials.
type Refresh struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

// Client implementes the client side of the OAuth flow.
type Client struct{ CacheLocation string }

// Exchange the code for an authentication token.
func (oauthClient *Client) Exchange(authCode string) Credentials {
	url := "https://api.freshbooks.com/auth/oauth/token"

	payload := Access{
		GrantType:    "authorization_code",
		ClientSecret: clientSecret,
		Code:         authCode,
		ClientID:     clientID,
		RedirectURI:  redirectURI,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("API-Version", "alpha")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Unexpected error when authenticating credentials (", resp.StatusCode, ")")
	}
	fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var credentials Credentials
	if err := json.Unmarshal(body, &credentials); err != nil {
		panic(err)
	}

	return credentials
}

// Refresh a stale authentication token for a new one.
func (oauthClient *Client) Refresh(credentials Credentials) Credentials {
	url := "https://api.freshbooks.com/auth/oauth/token"

	payload := Refresh{
		GrantType:    "refresh_token",
		RefreshToken: credentials.RefreshToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("API-Version", "alpha")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Unexpected error when refreshing credentials (", resp.StatusCode, ")")
	}
	fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var refreshed Credentials
	if err := json.Unmarshal(body, &refreshed); err != nil {
		panic(err)
	}

	return refreshed
}

// IsExpired determines if the token is still valid.
func (oauthClient *Client) IsExpired(credentials Credentials) bool {
	createdAt := time.Unix(int64(credentials.CreatedAt), 0)

	expiresAt := createdAt.Add(time.Second * time.Duration(int64(credentials.ExpiresIn)))

	expiredSecs := expiresAt.Sub(time.Now())
	return expiredSecs <= 0
}

func (oauthClient *Client) getCacheLocation() string {
	var location string
	if oauthClient.CacheLocation == "" {
		location = "~/.ttrack.creds.json"
	} else {
		location = oauthClient.CacheLocation
	}
	if strings.Contains(location, "~") {
		expanded, err := homedir.Expand(location)
		if err != nil {
			panic(err)
		}
		location = expanded
	}
	return location
}

// IsAuthenticated determines if the user is logged in. This does not
// actually verify with the service. It only checks to see if the
// credentials exist.
func (oauthClient *Client) IsAuthenticated() bool {
	_, err := oauthClient.FromCache()
	if err != nil && os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Fatal(err)
	}

	return true
}

// FromCache attempts to read existing oauth credentials from a cache.
func (oauthClient *Client) FromCache() (Credentials, error) {
	location := oauthClient.getCacheLocation()
	content, err := ioutil.ReadFile(location)
	if err != nil {
		return Credentials{}, err
	}
	credentials := Credentials{}
	err = json.Unmarshal(content, &credentials)
	return credentials, err
}

// Cache saves credentials to a local file.
func (oauthClient *Client) Cache(credentials Credentials) {
	location := oauthClient.getCacheLocation()
	data, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Writing credentials to:", location)
	if err := ioutil.WriteFile(location, data, 0644); err != nil {
		log.Fatal(err)
	}
}
