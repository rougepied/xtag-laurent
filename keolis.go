// Copyright © 2015 Arnaud Malivoir
// This work is free. You can redistribute it and/or modify it under the
// terms of the Do What The Fuck You Want To Public License, Version 2,
// as published by Sam Hocevar. See the COPYING file or http://www.wtfpl.net/
// for more details.

package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// OpenData represents the main struct return by the Keolis-Rennes API.
type OpenData struct {
	// Request is a copy of the request passed to the Keolis-Rennes API
	Request string  `xml:"request" json:"request"`
	Answer  *Answer `xml:"answer" json:"answer"`
}

// Answer is the answer contained in OpenData.
type Answer struct {
	Status Status `xml:"status" json:"status"`
	Data   *Data  `xml:"data" json:"data"`
}

// Status contains the status of the response.
//
//    0    Success
//    1    Invalid key
//    2    Invalid version
//    3    Invalid command
//    4    Empty key
//    5    Empty version
//    6    Empty command
//    8    Usage limit reached
//    98   Disabled
//    99   Maintenance
//    100+ Command error code
//
// see: http://data.keolis-rennes.com/fr/les-donnees/fonctionnement-de-lapi.html
type Status struct {
	Code    int    `xml:"code,attr" json:"code"`
	Message string `xml:"message,attr" json:"status"`
}

// Data of the Answer.
type Data struct {
	Stations      *Stations `xml:"station" json:"station"`
	LocalDateTime string    `xml:"localdatetime,attr" json:"localdatetime"`
	StopLine      *StopLine `xml:"stopline" json:"stopline"`
}

// StopLine of the Data
type StopLine struct {
	Stop       string       `xml:"stop" json:"stop"`
	Route      string       `xml:"route" json:"route"`
	Direction  int          `xml:"direction" json:"direction"`
	Departures []*Departure `xml:"departures>departure" json:"departures"`
}

// Departure represente the next bus / metro departure.
type Departure struct {
	Accurate  int    `xml:"accurate,attr" json:"accurate"`
	HeadSign  string `xml:"headsign,attr" json:"headsign"`
	Vehicle   int    `xml:"vehicle,attr" json:"vehicle"`
	Expected  string `xml:"expected,attr" json:"expected"`
	timeValue `xml:",innerxml"`
}

// Stopline of a Departure.
type Stopline struct {
	Stop       string     `xml:"stop"`
	Route      string     `xml:"route"`
	Direction  string     `xml:"direction"`
	Departures Departures `xml:"departures"`
}

// Departures is a list of Departure
type Departures []*Departure

func (d Departures) String() {
	for _, depart := range d {
		fmt.Printf("%#v\n", depart)
	}
}

// Stations is a list of Station
type Stations []*Station

// Station for a bus or a metro
type Station struct {
	ID             int     `xml:"id" json:"id"`
	Number         int     `xml:"number" json:"number"`
	Name           string  `xml:"name" json:"name"`
	State          int     `xml:"state" json:"state"`
	Latitude       float64 `xml:"latitude" json:"latitude"`
	Longitude      float64 `xml:"longitude" json:"longitude"`
	SlotsAvailable int     `xml:"slotsavailable" json:"slotsavailable"`
	BikesAvailable int     `xml:"bikesavailable" json:"bikesavailable"`
	Pos            int     `xml:"pos" json:"pos"`
	District       string  `xml:"district" json:"district"`
	LastUpdate     string  `xml:"lastupdate" json:"lastupdate"`
}

// timeValue is a custom type that contain a time.Time anonymous field.
// It satisfy the xml.Unmarshaler interface so that it can be used to map
// responses from Keolis Rennes API
// See: http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields/25015260#25015260
type timeValue struct {
	time.Time
}

func (t *timeValue) UnmarshalXMLAttr(attr xml.Attr) error {
	fmt.Printf("Parsing attribute '%s', with value '%s'\n", attr.Name.Local, attr.Value)
	const shortForm = "2006-Jan-02"
	parse, err := time.Parse(shortForm, "2013-Feb-03")
	if err != nil {
		return nil
	}
	*t = timeValue{parse}
	return nil
}

// UnmarshalXML unmarshal the timeValue according to the specific format
func (t *timeValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse("2006-01-02T15:04:05-07:00", v)
	if err != nil {
		return err
	}
	*t = timeValue{parse}
	return nil
}

// Schedule is the content of Response
type Schedule struct {
	Time timeValue `xml:"time" json:"time"`
	Line string    `xml:"line" json:"line"`
}

// Schedules is a list of Schedule
type Schedules []*Schedule

// Append a Schedule to the Schedules
func (s *Schedules) Append(item *Schedule) {
	*s = append(*s, item)
}

// Response is the response to "/api/3.0"
type Response struct {
	Schedules *Schedules `xml:"schedules" json:"schedules"`
	Error     string     `xml:"error" json:"error"`
}

// NewResponse create a new Response
func NewResponse() *Response {
	return &Response{Schedules: &Schedules{}}
}

// handleAPI3 handle response to "/api/3.0"
func handleAPI3(key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parsing body

		// vars := mux.Vars(r)
		stopID := r.FormValue("stop")
		routeID := r.FormValue("route")
		directionID := r.FormValue("direction")

		// getting data from Keolis
		ret, err := getBusNextDepartures3(r, key, stopID, routeID, directionID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			debugf(r, "%s", err.Error())
			return
		}

		// responding
		headers := w.Header()
		headers["Content-Type"] = []string{"application/json"}

		msg, err := json.Marshal(ret)
		if err != nil {
			http.Error(w, err.Error(), 500)
			debugf(r, "%s", err.Error())
			return
		}

		//		debugf(r, "réponse : %s", msg)
		fmt.Fprintf(w, "%s", msg)
	})
}

// getBusNextDepartures3 get the next departures for handleAPI3
func getBusNextDepartures3(r *http.Request, key, stopID, routeID, directionID string) (*Response, error) {
	var URL *url.URL
	URL, err := url.Parse("http://data.keolis-rennes.com/xml/")
	if err != nil {
		debugf(r, "%s", err.Error())
		return nil, err
	}
	parameters := url.Values{}
	parameters.Add("key", key)
	parameters.Add("cmd", "getbusnextdepartures")
	parameters.Add("version", "2.2")
	parameters.Add("param[mode]", "stopline")
	parameters.Add("param[stop][]", stopID)
	parameters.Add("param[route][]", routeID)
	parameters.Add("param[direction][]", directionID)
	URL.RawQuery = parameters.Encode()

	//fmt.Println(URL.String())

	resp, err := get(r, URL.String())
	if err != nil {
		debugf(r, "%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debugf(r, "%s", err.Error())
		return nil, err
	}

	o, err := unmarshalResponse(body)
	if err != nil {
		debugf(r, "%s", err.Error())
		return nil, err
	}

	answer := o.Answer
	if answer == nil {
		debugf(r, "Response %v\n", o)
		return nil, err
	}

	response := NewResponse()
	if answer.Status.Code != 0 {
		response.Error = answer.Status.Message
		return response, err
	}

	if answer.Data == nil || answer.Data.StopLine == nil {
		debugf(r, "Response %v\n", o)
		return nil, err
	}

	for _, d := range answer.Data.StopLine.Departures {
		response.Schedules.Append(&Schedule{
			Line: answer.Data.StopLine.Route,
			Time: d.timeValue,
		})
	}

	return response, err
}

// unmarshalResponse convert []byte to OpenData
func unmarshalResponse(data []byte) (OpenData, error) {
	decoder := xml.NewDecoder(bytes.NewBuffer(data))
	decoder.Strict = false
	var o OpenData
	err := decoder.Decode(&o)

	return o, err
}
