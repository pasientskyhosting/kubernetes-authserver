package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Index",
		"HEAD",
		"/",
		Index,
	},
	Route{
		"Auth",
		"POST",
		"/auth",
		Auth,
	},
	Route{
		"Healthz",
		"GET",
		"/healthz",
		Healthz,
	},
	Route{
		"Healthz",
		"HEAD",
		"/healthz",
		Healthz,
	},
}
