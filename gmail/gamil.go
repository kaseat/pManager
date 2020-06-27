package gmail

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"

	"github.com/kaseat/pManager/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type client struct {
	config *oauth2.Config
	ctx    context.Context
}

var cl *client

// Client is Gmail client
type Client interface {
	// Returns url handled by GMail to identify user
	GetAuthURL(login string) (string, error)
	// Handles response from gmail request
	HandleResponse(state string, code string) error
	// Gets gmail client for given login
	GetServiceForUser(login string) (*gmail.Service, error)
}

// GetClient gets Gmail client
func GetClient() Client {
	if cl == nil {
		cfg, _ := getServiceFromFile()
		cl = &client{
			config: cfg,
			ctx:    context.Background(),
		}
	}
	return *cl
}

// GetAuthUrl returns url handled by GMail to identify user
func (c client) GetAuthURL(login string) (string, error) {
	stateRaw := make([]byte, 18)
	if _, err := rand.Read(stateRaw); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(stateRaw)

	s := storage.GetStorage()
	err := s.AddUserState(login, state)
	if err != nil {
		return "", err
	}

	authURL := c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return authURL, nil
}

func (c client) HandleResponse(state string, code string) error {

	tok, err := c.config.Exchange(c.ctx, code)
	if err != nil {
		return err
	}

	s := storage.GetStorage()
	err = s.AddUserToken(state, tok)
	if err != nil {
		return err
	}
	return nil
}

func getServiceFromFile() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
}

func (c client) GetServiceForUser(login string) (*gmail.Service, error) {
	s := storage.GetStorage()
	tok, err := s.GetUserToken(login)
	if err != nil {
		return nil, err
	}

	return gmail.New(c.config.Client(c.ctx, &tok))
}
