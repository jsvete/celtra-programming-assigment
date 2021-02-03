// Package rest contains handler code for REST API calls
package rest

import (
	"celtra-programming-assigment/cmd/tracker/persistence"
	"celtra-programming-assigment/pkg/pubsub"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// CreateRouter returns a router with registered handlers
func CreateRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/:accountId", handleGet)
	router.POST("/", handlePost)
	router.PUT("/:accountId", handlePut)

	return router
}

// handleGet function handles GET requests.
//
// It returns a JSON representation of an account matching the accountID (e.g. GET BASE_URL/{accountID}).
func handleGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	accountID, err := parseAccountID(params)
	if err != nil {
		log.Error().Msgf("invalid accountId value: %v", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	account, err := persistence.DB.GetAccount(accountID)
	if err != nil {
		log.Error().Msgf("getting account %d from database: %v", accountID, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, err := json.Marshal(account)
	if err != nil {
		log.Error().Msgf("serializing account %d to JSON: %v", accountID, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Error().Msgf("writing body for account %d: %v", accountID, err)
	}
}

// handlePost function handles POST requests.
//
// It creates a new account and returns a Location header with a relative URL where the new account can be accesed from.
//
// The function accepts JSON payload in the following format: {"name":"ACCOUNT_NAME", "isActive": true/false}
func handlePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("reading body: %v", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	bodyStruct := struct {
		Name     string
		IsActive bool
	}{}

	if err := json.Unmarshal(bodyBytes, &bodyStruct); err != nil {
		log.Error().Msgf("invalid JSON format in the body: %v", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	account, err := persistence.DB.CreateAccount(bodyStruct.Name, bodyStruct.IsActive)
	if err != nil {
		log.Error().Msgf("creating new account: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Location", fmt.Sprintf("/%d", account.ID))
	w.WriteHeader(http.StatusCreated)
}

// handlePut function handles PUT requests.
//
// It is used to receive events for a specific account (e.g. PUT BASE_URL/{accountID}?data="ACCOUNT_DATA")
func handlePut(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	accountID, err := parseAccountID(params)
	if err != nil {
		log.Error().Msgf("invalid accountId value: %v", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	active, err := persistence.DB.IsActiveAccount(accountID)
	if err != nil {
		log.Error().Msgf("checking if accountID %d is active: %v", accountID, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if !active {
		log.Error().Msgf("accoundID %d is not active: %v", accountID, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	data := r.URL.Query().Get("data")
	if data == "" {
		log.Error().Msgf("missing data value for accoundID %d: %v", accountID, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	go func() {
		if err := pubsub.Bus.Publish(accountID, data); err != nil {
			log.Error().Msgf("publishing event for accoundID %d: %v", accountID, err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

// parseAccountID is a helper function to parse account ID from the request context.
func parseAccountID(params httprouter.Params) (int, error) {
	accountIDParam := params.ByName("accountId")

	if accountIDParam == "" {
		return -1, errors.New("missing accountId parameter")
	}

	accountID, err := strconv.Atoi(accountIDParam)
	if err != nil {
		return -1, fmt.Errorf("invalid accountId %s: %v", accountIDParam, err)
	}

	return accountID, nil
}
