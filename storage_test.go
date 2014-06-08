package osin

import (
	"errors"
	"strconv"
)

type TestingStorage struct {
	clients   map[string]*Client
	authorize map[string]*AuthorizeData
	access    map[string]*AccessData
	refresh   map[string]string
}

func NewTestingStorage() *TestingStorage {
	r := &TestingStorage{
		clients:   make(map[string]*Client),
		authorize: make(map[string]*AuthorizeData),
		access:    make(map[string]*AccessData),
		refresh:   make(map[string]string),
	}

	r.clients["1234"] = &Client{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth",
	}

	return r
}

func (s *TestingStorage) GetClient(id string) (*Client, error) {
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, errors.New("Client not found")
}

func (s *TestingStorage) SetClient(id string, client *Client) error {
	s.clients[id] = client
	return nil
}

func (s *TestingStorage) SaveAuthorize(data *AuthorizeData) error {
	s.authorize[data.Code] = data
	return nil
}

func (s *TestingStorage) LoadAuthorize(code string) (*AuthorizeData, error) {
	if d, ok := s.authorize[code]; ok {
		return d, nil
	}
	return nil, errors.New("Authorize not found")
}

func (s *TestingStorage) RemoveAuthorize(code string) error {
	delete(s.authorize, code)
	return nil
}

func (s *TestingStorage) SaveAccess(data *AccessData) error {
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

func (s *TestingStorage) LoadAccess(code string) (*AccessData, error) {
	if d, ok := s.access[code]; ok {
		return d, nil
	}
	return nil, errors.New("Access not found")
}

func (s *TestingStorage) RemoveAccess(code string) error {
	delete(s.access, code)
	return nil
}

func (s *TestingStorage) LoadRefresh(code string) (*AccessData, error) {
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, errors.New("Refresh not found")
}

func (s *TestingStorage) RemoveRefresh(code string) error {
	delete(s.refresh, code)
	return nil
}

// Predictable testing token generation

type TestingAuthorizeTokenGen struct {
	counter int64
}

func (a *TestingAuthorizeTokenGen) GenerateAuthorizeToken(data *AuthorizeData) (ret string, err error) {
	a.counter++
	return strconv.FormatInt(a.counter, 10), nil
}

type TestingAccessTokenGen struct {
	acounter int64
	rcounter int64
}

func (a *TestingAccessTokenGen) GenerateAccessToken(data *AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	a.acounter++
	accesstoken = strconv.FormatInt(a.acounter, 10)

	if generaterefresh {
		a.rcounter++
		refreshtoken = "r" + strconv.FormatInt(a.rcounter, 10)
	}
	return
}