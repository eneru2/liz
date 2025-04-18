package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// TODO:
// - create a route for each folder containing a +page.templ
// - check if the folder has a layout.templ
// make a tree structure so if a route has a +layout.templ
// the childs of the route will inherit the layout
// or instead if a route has a layout.templ
// use a function that creates each child route with the layout

// NEW TODO
// load the root page
// recursion through routes
// how:
// when checking a folder check if inside that folder theres other folders
// if so

// NEW TODO
// experiment with creating liz files
// this files only contain imports
// and templates 
// and then dump the liz files
// into .templ files to avoid having
// to name each package, autoname with
// the folder name of the route

type Route struct {
	Name string
	HasHead bool
	HasBody bool
}


func main() {
	// first scan src/routes to generate a route for each folder
	// containing a +page.templ
	os.RemoveAll(".liz/files")
	os.Mkdir(".liz/files", 0755)
	files, _ := os.ReadDir("src/routes")
	
	// only checking if theres a base path "/"`
	routes := checkForRootIndex([]Route{})
	fmt.Println(routes)

	routes = checkRouteForFolders(0, files, routes)
	fmt.Println(routes)

	for _, route := range routes {
		createComponent(route)
	}
	createHandlers(routes)
}

func createHandlers(routes []Route) {
	var routeHandlers string
	for _, route := range routes {
		if route.Name != "root" {
			routeHandlers += `r.Get("/`+route.Name+`",`
			if strings.Contains(route.Name, "/") {
				route.Name = strings.ReplaceAll(route.Name, "/", "_")
			}
			routeHandlers += `templ.Handler(Route_`+route.Name+`()).ServeHTTP)` + "\n"
		} else {
		routeHandlers += `r.Get("/", templ.Handler(Route_`+route.Name+`()).ServeHTTP)` + "\n"
		}
	}
	var template = `package files

import (
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux) {
	`+routeHandlers+`
}
`
	// os.WriteFile("liz/routes/handlers/exported_routes.go", []byte(template), 0644)
	os.WriteFile(".liz/files/exported_routes.go", []byte(template), 0644)
}

func createComponent(route Route) {
	if !route.HasHead && !route.HasBody {
		return
	}
	var Page string
	var selfImport string
	if route.Name != "root" {
		if strings.Contains(route.Name, "/") {
			no_slash := strings.ReplaceAll(route.Name, "/", "_")
			selfImport = ``+no_slash+` "liz/src/routes/`+route.Name+`"`
		} else {
			selfImport = `"liz/src/routes/`+route.Name+`"`
		}
	} else {
		selfImport = `root "liz/src/routes"`
	}

	if strings.Contains(route.Name, "/") {
		route.Name = strings.ReplaceAll(route.Name, "/", "_")
	}

	if route.HasBody && !route.HasHead {
		Page = `var Page_`+route.Name+` = app.Page{
			HeadContents: nil,
			BodyContents: `+route.Name+`.Body(),
		}`
	}
	if !route.HasBody && route.HasHead {
		Page = `var Page_`+route.Name+` = app.Page{
			HeadContents: `+route.Name+`.Head(),
			BodyContents: nil,
		}`
	}
	if route.HasHead && route.HasBody {
		Page = `var Page_`+route.Name+` = app.Page{
			HeadContents: `+route.Name+`.Head(),
			BodyContents: `+route.Name+`.Body(),
		}`
	}
		var template = `package files

import (
	"github.com/a-h/templ"
	"liz/src"
	`+selfImport+`
)

`+Page+`

func Route_`+route.Name+`() templ.Component {
	return app.DefaultLayout(Page_`+route.Name+`)
}`

	os.WriteFile(".liz/files/"+route.Name+".go", []byte(template), 0644)
}
func checkForRootIndex(routes []Route) []Route {
	path := "src/routes/page.templ"
	_, err := os.Stat(path)
	if err != nil {
		return routes
	}

	route := Route{
		Name: "root",
		HasHead: routeHasHead(path),
		HasBody: routeHasBody(path),
	}

	newRoutes := append(routes, route)

	return newRoutes
}

func checkSubroutes(index int, files []fs.DirEntry, routes []Route, parentsNames string) []Route {
	if index < len(files) {
		if files[index].IsDir() {
			path := "src/routes/" + parentsNames + "/" + files[index].Name() + "/page.templ"
			_, err := os.Stat(path)
			if err != nil {
				// this is probably bad, but cant come up 
				// with a better solution now
				panic(err)
			}

			route := Route{
				Name: parentsNames + "/" + files[index].Name(),
				HasHead: routeHasHead(path),
				HasBody: routeHasBody(path),
			}
			// files, _ := os.ReadDir("src/routes/"+parentsNames+"/"+files[index].Name())

			routes = append(routes, route)
		}
		checkSubroutes(index+1, files, routes, parentsNames)
	}
	return routes
}

func checkRouteForFolders(index int, files []fs.DirEntry, routes []Route) []Route {
	if index < len(files) {
		// check if the folder has a +page.templ
		if files[index].IsDir() {
			path := "src/routes/" + files[index].Name() + "/page.templ"
			_, err := os.Stat(path)
			if err != nil {
				// this is probably bad, but cant come up 
				// with a better solution now
				panic(err)
			}

			route := Route{
				Name: files[index].Name(),
				HasHead: routeHasHead(path),
				HasBody: routeHasBody(path),
			}
			fmt.Println(route.Name)
			files, _ := os.ReadDir("src/routes/"+files[index].Name())
			routes = append(routes, route)

			routes = checkSubroutes(0, files, routes, route.Name)
			
		}
		checkRouteForFolders(index+1, files, routes)
	}
	fmt.Println(routes)
	return routes
}

func routeHasHead(path string) bool {
	headRegex := `templ Head(`
	file, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	if strings.Contains(string(file), headRegex) {
		return true
	}
	
	return false
}

func routeHasBody(path string) bool {
	bodyRegex := `templ Body(`
	file, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	if strings.Contains(string(file), bodyRegex) {
		return true
	}
	
	return false
}