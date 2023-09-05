package server

import (
	"fmt"
	"net/http"
	"sync"
)

type App struct {
	sync sync.Once
	mux  http.ServeMux
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.sync.Do(func() {
		app.mux = http.ServeMux{}
		app.mux.HandleFunc("/", app.Home)
		app.mux.HandleFunc("/alert", app.Alert)
	})
	app.mux.ServeHTTP(w, r)
}

func (app *App) Alert(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
	<html>
	  <head><title>Alert Page</title></head>
	  <body>
	    <div class="alert alert-primary" role="alert">
		  Alert!
		  Something is Wrong!
		  Alert!
		</div>
		<h1>Alert Page</h1>
	  </body>
	</html>`)
}

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
	<html>
	  <head><title>Home Page</title></head>
	  <body>
	    <h1>Hi!</h1>
		<p>This is the home page...</p>
	  </body>
	</html>`)
}
