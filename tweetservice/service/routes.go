package service

import "net/http"

// Defines a single route, e.g. a human readable name, HTTP method and the
// pattern the function that will execute when the route is called.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Defines the type Routes which is just an array (slice) of Route structs.
type Routes []Route

// Initialize our routes
var routes = Routes{

	Route{
		"GetTweets", // Name
		"GET",        // HTTP method
		"/tweets/{accountId}", // Route pattern
		GetTweets,
	},
	Route{
		"GetTweetsPaginated",
		"GET",
		"/get/{offset}",
		GetTweetsPaginated,
	},
	Route{
		"AddTweet",
		"GET",
		"/add",
		AddTweet,
	},
	Route{
		"HealthCheck",
		"GET",
		"/health",
		HealthCheck,
	},
	Route{
		"Testability",
		"GET",
		"/testability/healthy/{state}",
		SetHealthyState,
	},
}
