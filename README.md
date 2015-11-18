# xtag-laurent

xtag-Laurent is a small web app showing the next bus departures at Saint-Laurent station, Rennes, France.

My goal is to have fun with:
- [Go](http://golang.org/) and [Google App Engine](https://cloud.google.com/appengine/)
- [Web Components](http://webcomponents.org/) written in [ES2015](https://babeljs.io/docs/learn-es2015/) and [X-Tag](x-tag.github.io)
- [Brunch.io](http://brunch.io/) to build statics
- Keolis Rennes API and open data

# Getting the code

## Witch git an go

	git clone https://github.com/rougepied/saint-laurent.git

## Pure Go App

	go get github.com/rougepied/saint-laurent

## Appengine

	goapp get github.com/rougepied/saint-laurent

_Note_: The `goapp` tool is part of the [Google App Engine SDK for Go](https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go).

# Front end

Assuming that you have node (with npm) bower and [Brunch.io](http://brunch.io/) installed.

	npm install
	bower install
	brunch b -P
