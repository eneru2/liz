package files

import (
	"github.com/a-h/templ"
	"liz/src"
	root "liz/src/routes"
)

var Page_root = app.Page{
			HeadContents: nil,
			BodyContents: root.Body(),
		}

func Route_root() templ.Component {
	return app.DefaultLayout(Page_root)
}