// Copyright © 2015 Arnaud Malivoir
// This work is free. You can redistribute it and/or modify it under the
// terms of the Do What The Fuck You Want To Public License, Version 2,
// as published by Sam Hocevar. See the COPYING file or http://www.wtfpl.net/
// for more details.

// +build appengine

// A Google App Engine application providing a web UI showing next bus departures.
package main

import (
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.Handle("/api/3.0", handleAPI3(keyDefault))
}

func debugf(r *http.Request, format string, args ...interface{}) {
	c := appengine.NewContext(r)
	c.Debugf(format, args)
}

func get(r *http.Request, url string) (resp *http.Response, err error) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	return client.Get(url)
}
