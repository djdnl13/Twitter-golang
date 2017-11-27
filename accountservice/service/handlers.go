package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/djdnl13/twitter/accountservice/dbclient"
	"github.com/djdnl13/twitter/accountservice/model"
	"github.com/djdnl13/twitter/common/messaging"
	"github.com/djdnl13/twitter/common/util"
	"github.com/gorilla/mux"
)

var DBClient dbclient.IBoltClient
var MessagingClient messaging.IMessagingClient
var isHealthy = true

var client = &http.Client{}

var LOGGER = logrus.Logger{}

func init() {
	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	client.Transport = transport
	LOGGER.Infof("Successfully initialized transport")
}

func GetAccount(w http.ResponseWriter, r *http.Request) {

	// Read the 'accountId' path parameter from the mux map
	var accountId = mux.Vars(r)["accountId"]

	// Read the account struct BoltDB
	account, err := DBClient.QueryAccount(accountId)
	account.ServedBy = util.GetIP()

	// If err, return a 404
	if err != nil {
		logrus.Errorf("Some error occured serving " + accountId + ": " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	notifyVIP(account) // Send VIP notification concurrently

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(account)
	writeJsonResponse(w, http.StatusOK, data)
}

// If our hard-coded "VIP" account, spawn a goroutine to send a message.
func notifyVIP(account model.Account) {
}

func SetHealthyState(w http.ResponseWriter, r *http.Request) {

	// Read the 'state' path parameter from the mux map and convert to a bool
	var state, err = strconv.ParseBool(mux.Vars(r)["state"])

	// If we couldn't parse the state param, return a HTTP 400
	if err != nil {
		logrus.Errorln("Invalid request to SetHealthyState, allowed values are true or false")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Otherwise, mutate the package scoped "isHealthy" variable.
	isHealthy = state
	w.WriteHeader(http.StatusOK)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Since we're here, we already know that HTTP service is up. Let's just check the state of the boltdb connection
	dbUp := DBClient.Check()
	if dbUp && isHealthy { // NEW condition here!
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJsonResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Database unaccessible"})
		writeJsonResponse(w, http.StatusServiceUnavailable, data)
	}
}

func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
