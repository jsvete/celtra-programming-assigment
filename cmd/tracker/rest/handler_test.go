// Package rest contains handler code for REST API calls
package rest

import (
	"bytes"
	"celtra-programming-assigment/pkg/dto"
	"celtra-programming-assigment/pkg/persistence"
	"celtra-programming-assigment/pkg/pubsub"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockedDB implements persistence.Database interface and exposes
// functions that can be used to mock the database response.
type mockedDB struct {
	FnIsActiveAccount func(ID int) (bool, error)
	FnCreateAccount   func(name string, isActive bool) (*dto.Account, error)
	FnGetAccount      func(ID int) (*dto.Account, error)
}

func (m *mockedDB) IsActiveAccount(ID int) (bool, error) {
	if m.FnIsActiveAccount == nil {
		return false, errorNotImplemented
	}

	return m.FnIsActiveAccount(ID)
}

func (m *mockedDB) CreateAccount(name string, isActive bool) (*dto.Account, error) {
	if m.FnCreateAccount == nil {
		return nil, errorNotImplemented
	}

	return m.FnCreateAccount(name, isActive)
}

func (m *mockedDB) GetAccount(ID int) (*dto.Account, error) {
	if m.FnGetAccount == nil {
		return nil, errorNotImplemented
	}

	return m.FnGetAccount(ID)
}

// mockedBus implements pubsub.PubSub interface and exposes
// functions that can be used to mock the publish/subscribe calls.
type mockedBus struct {
	FnPublish   func(accountID int, data string) error
	FnSubscribe func() chan *pubsub.Event
}

func (b *mockedBus) Publish(accountID int, data string) error {
	if b.FnPublish == nil {
		return errorNotImplemented
	}

	return b.FnPublish(accountID, data)
}

func (b *mockedBus) Subscribe() chan *pubsub.Event {
	if b.FnSubscribe == nil {
		return nil
	}

	return b.FnSubscribe()
}

var (
	server              *httptest.Server
	errorNotImplemented = errors.New("not implemented")
	fakeDB              *mockedDB
	fakeBus             *mockedBus
	accounts            = map[int]*dto.Account{}
)

func TestMain(m *testing.M) {
	// setup a testing server with registered routes
	server = httptest.NewServer(CreateRouter())
	defer server.Close()

	// database mock
	fakeDB = &mockedDB{}
	persistence.DB = fakeDB

	fmt.Printf("db: %+v\n", persistence.DB)

	// pubsub mock
	fakeBus = &mockedBus{}
	pubsub.Bus = fakeBus

	fmt.Printf("pubsub: %+v\n", pubsub.Bus)

	m.Run()
}

func Test_Post(t *testing.T) {
	fakeDB.FnCreateAccount = func(name string, isActive bool) (*dto.Account, error) {
		account := &dto.Account{
			ID:       1,
			Name:     name,
			IsActive: isActive,
		}
		accounts[account.ID] = account

		fmt.Printf("%+v\n", account)
		fmt.Printf("%+v\n", accounts)
		return account, nil
	}

	bodyStruct := struct {
		Name     string
		IsActive bool
	}{
		Name:     "test account",
		IsActive: true,
	}

	body, err := json.Marshal(&bodyStruct)
	if err != nil {
		t.Fatalf("marshaling body: %v", err)
	}

	req, err := http.NewRequest("POST", server.URL, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	locationHeader := resp.Header.Get("Location")

	if locationHeader != "/1" {
		t.Fatalf("expected %s but got %s", "/1", locationHeader)
	}
}

func Test_PostBadDatabase(t *testing.T) {
	fakeDB.FnCreateAccount = func(name string, isActive bool) (*dto.Account, error) {
		return nil, errors.New("bad database")
	}

	bodyStruct := struct {
		Name     string
		IsActive bool
	}{
		Name:     "test account",
		IsActive: true,
	}

	body, err := json.Marshal(&bodyStruct)
	if err != nil {
		t.Fatalf("marshaling body: %v", err)
	}

	req, err := http.NewRequest("POST", server.URL, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected %d but got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func Test_PostBadBody(t *testing.T) {
	body, err := json.Marshal("bad body")
	if err != nil {
		t.Fatalf("marshaling body: %v", err)
	}

	req, err := http.NewRequest("POST", server.URL, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func Test_Get(t *testing.T) {
	fakeDB.FnGetAccount = func(ID int) (*dto.Account, error) {
		return accounts[ID], nil
	}

	url := server.URL + "/1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d but got %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected %s but got %s", "application/json", contentType)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	account := &dto.Account{}
	if err := json.Unmarshal(body, account); err != nil {
		t.Fatalf("unmarshaling body: %v", err)
	}

	if account.ID != 1 {
		t.Fatalf("expected %d but got %d", 1, account.ID)
	}

	if account.Name != "test account" {
		t.Fatalf("expected %s but got %s", "test account", account.Name)
	}

	if account.IsActive != true {
		t.Fatalf("account is not active")
	}
}

func Test_GetInvalidParam(t *testing.T) {
	url := server.URL + "/asd"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func Test_GetNoAccount(t *testing.T) {
	fakeDB.FnGetAccount = func(ID int) (*dto.Account, error) {
		return nil, errors.New("no account")
	}

	url := server.URL + "/666"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected %d but got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func Test_Put(t *testing.T) {
	fakeDB.FnIsActiveAccount = func(ID int) (bool, error) {
		account, ok := accounts[ID]
		if !ok {
			return false, errors.New("no account")
		}

		return account.IsActive, nil
	}

	fakeBus.FnPublish = func(accountID int, data string) error {
		return nil
	}

	urlWithData := server.URL + "/1?data=testdata"

	req, err := http.NewRequest("PUT", urlWithData, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d but got %d", http.StatusAccepted, resp.StatusCode)
	}
}

func Test_PutNotActive(t *testing.T) {
	fakeDB.FnIsActiveAccount = func(ID int) (bool, error) {
		return false, nil
	}

	fakeBus.FnPublish = func(accountID int, data string) error {
		return nil
	}

	urlWithData := server.URL + "/1?data=testdata"

	req, err := http.NewRequest("PUT", urlWithData, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func Test_PutNoData(t *testing.T) {
	fakeDB.FnIsActiveAccount = func(ID int) (bool, error) {
		account, ok := accounts[ID]
		if !ok {
			return false, errors.New("no account")
		}

		return account.IsActive, nil
	}

	fakeBus.FnPublish = func(accountID int, data string) error {
		return nil
	}

	url := server.URL + "/1"

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func Test_PutInvalidParam(t *testing.T) {
	url := server.URL + "/asd"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
