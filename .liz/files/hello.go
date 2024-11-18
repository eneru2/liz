package files

import (
	"github.com/a-h/templ"
	"liz/src"
	"liz/src/routes/hello"
)

var Page_hello = app.Page{
			HeadContents: nil,
			BodyContents: hello.Body(),
		}

func Route_hello() templ.Component {
	return app.DefaultLayout(Page_hello)
}