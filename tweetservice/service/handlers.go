package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/djdnl13/twitter/tweetservice/dbclient"
//	"github.com/djdnl13/twitter/tweetservice/model"
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

func AddTweet(w http.ResponseWriter, r *http.Request) {
	LOGGER.Infof("Successfully resolved add")
	r.ParseForm()

	var accountId string = r.Form.Get("accountId")
	var text string = r.Form.Get("text")
	var likesCount string = "0"

	result, err := DBClient.AddTweet(accountId, text, likesCount)

	// If err, return a 404
	if err != nil {
		logrus.Errorf("Some error occured saving tweet for account " + accountId + ": " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(result)

	writeJsonResponse(w, http.StatusOK, data)
}

func GetTweetsPaginated(w http.ResponseWriter, r *http.Request) {

	var offset = mux.Vars(r)["offset"]

	// Read the account struct BoltDB
	tweets, err := DBClient.QueryTweetsOffset(offset)

	//tweet.ServedBy = util.GetIP()
	for k, _ := range tweets {
		tweets[k].ServedBy = util.GetIP()
	}

	// If err, return a 404
	if err != nil {
		logrus.Errorf("Some error occured serving tweets offset" + offset + ": " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(tweets)
	writeJsonResponse(w, http.StatusOK, data)


}

func GetTweets(w http.ResponseWriter, r *http.Request) {

	// Read the 'accountId' path parameter from the mux map
	var accountId = mux.Vars(r)["accountId"]

	// Read the account struct BoltDB
	tweets, err := DBClient.QueryTweets(accountId)
	//tweet.ServedBy = util.GetIP()
	for k, v := range tweets {
		tweets[k].ServedBy = util.GetIP()
		v.ServedBy = util.GetIP()
	}

	// If err, return a 404
	if err != nil {
		logrus.Errorf("Some error occured serving " + accountId + ": " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(tweets)
	writeJsonResponse(w, http.StatusOK, data)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
