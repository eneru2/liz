package files

import (
	"github.com/a-h/templ"
	"liz/src"
	hello_hey "liz/src/routes/hello/hey"
)

var Page_hello_hey = app.Page{
			HeadContents: nil,
			BodyContents: hello_hey.Body(),
		}

func Route_hello_hey() templ.Component {
	return app.DefaultLayout(Page_hello_hey)
}