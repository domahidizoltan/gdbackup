package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type Oauth2Client struct {
	config *oauth2.Config
}

func NewOauth2Client() *http.Client {
	creds, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(creds, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := Oauth2Client{
		config: config,
	}
	return client.getClient()
}

// Retrieve a token, saves the token, then returns the generated client.
func (client Oauth2Client) getClient() *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := client.tokenFromFile(tokFile)
	if err != nil {
		tok = client.getTokenFromWeb()
		client.saveToken(tokFile, tok)
	}
	return client.config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func (client Oauth2Client) getTokenFromWeb() *oauth2.Token {
	authURL := client.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := client.config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func (client Oauth2Client) tokenFromFile(tokenFile string) (*oauth2.Token, error) {
	file, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (client Oauth2Client) saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
}
