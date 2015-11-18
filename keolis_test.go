// Copyright © 2015 Arnaud Malivoir
// This work is free. You can redistribute it and/or modify it under the
// terms of the Do What The Fuck You Want To Public License, Version 2,
// as published by Sam Hocevar. See the COPYING file or http://www.wtfpl.net/
// for more details.

package main

import (
	"fmt"
	"testing"
)

var fakedResponse = []byte(`
<opendata>
<request>http://data.keolis-rennes.com/xml/?cmd=getbusnextdepartures&version=2.2&key=xxxxxxxxx&param[mode]=stopline&param[stop][]=1372&param[route][]=0009&param[direction][]=0</request>
<answer>
	<status code="0" message="OK"/>
	<data localdatetime="2015-05-21T13:35:54+02:00">
		<stopline>
			<stop>1372</stop>
			<route>0009</route>
			<direction>0</direction>
			<departures>
				<departure accurate="1" headsign="Cleunay" vehicle="268447282" expected="2015-05-21T13:43:00+02:00">2015-05-21T13:43:00+02:00</departure>
				<departure accurate="1" headsign="Cleunay" vehicle="268447119" expected="2015-05-21T13:52:00+02:00">2015-05-21T13:52:00+02:00</departure>
				<departure accurate="1" headsign="Cleunay" vehicle="268447222" expected="2015-05-21T14:02:00+02:00">2015-05-21T14:02:00+02:00</departure>
				<departure accurate="0" headsign="Cleunay" vehicle="268447183" expected="2015-05-21T14:12:00+02:00">2015-05-21T14:12:00+02:00</departure>
				<departure accurate="0" headsign="Cleunay" vehicle="268447098" expected="2015-05-21T14:21:00+02:00">2015-05-21T14:21:00+02:00</departure>
				<departure accurate="0" headsign="Cleunay" vehicle="268447203" expected="2015-05-21T14:29:00+02:00">2015-05-21T14:29:00+02:00</departure>
			</departures>
		</stopline>
	</data>
</answer>
</opendata>`)

func TestUnmarshalResponse(t *testing.T) {

	o, err := unmarshalResponse(fakedResponse)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", o.Answer.Data.StopLine.Departures)

	if o.Answer.Status.Code != 0 {
		t.Errorf("Status code = %d, expected %d", o.Answer.Status.Code, 0)
	}

}
