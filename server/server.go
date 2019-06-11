package server

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
	"shorty-url/db"
	"shorty-url/http_helpers"
)

type ShortyServer struct {
	router   *mux.Router
	database db.Database
}

func CreateServer(database db.Database) ShortyServer {
	s := ShortyServer{}
	s.database = database
	return s
}

func (a *ShortyServer) Init() {
	glog.Info("initializing database")
	a.database.Init()
	glog.Info("registering server paths")
	a.router = mux.NewRouter()
	a.router.HandleFunc("/", a.Home).Methods("GET")
	a.router.HandleFunc("/register", a.Register).Methods("POST")
	a.router.HandleFunc("/{p}", a.Redirect).Methods("GET")
}

//Run the app
func (a *ShortyServer) Run(port string) {
	glog.Infof("serving on 127.0.0.1%s", port)
	glog.Fatal(http.ListenAndServe(port, a.router))
}

func (a *ShortyServer) Home(w http.ResponseWriter, r *http.Request) {
	glog.Info("Home Page")
	http_helpers.PrintText(`
	<html lang="en">
	<head><title>Simple Shorty</title></head>
	<body>
	<form action="register" method="post">
	<input name="to" placeholder="destination URL"/>
	<input name="from" placeholder="shortened URL"/>
	<button type="submit" value="Submit"/>
	</form>
	</body>
	</html>
	`, w)
}

func (a *ShortyServer) Register(w http.ResponseWriter, r *http.Request) {

	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		glog.Errorf("ParseForm() err: %v", err)
		http_helpers.PrintError(http.StatusBadRequest, "Error parsing request", w)
		return
	}
	from := r.FormValue("from")
	to := r.FormValue("to")

	glog.Infof("received request to map %s to %s", from, to)

	if !http_helpers.IsValidURL(to) {
		http_helpers.PrintError(http.StatusBadRequest, "Invalid URL", w)
		return
	}
	if _, exists := a.database.Get(from); exists {
		glog.Infof("%s already exists", from)
		http_helpers.PrintError(http.StatusBadRequest, "URL Already Mapped", w)
		return
	}
	if err := a.database.AddMapping(from, to); err != nil {
		glog.Errorf("Encountered following error when adding mapping from: %s to: %s with error: %v", from, to, err)
		http_helpers.PrintError(http.StatusInternalServerError, "Internal Error", w)
		return
	}
	glog.Infof("mapped %s to %s", from, to)
	// Trusting that valid URLs won't have injection potential
	http_helpers.PrintText(fmt.Sprintf("Successfully registerd %s to direct to %s", from, to), w)
}

func (a *ShortyServer) Redirect(w http.ResponseWriter, r *http.Request) {
	from := mux.Vars(r)["p"]

	if redirectUrl, ok := a.database.Get(from); ok {
		http_helpers.Redirect(redirectUrl, w, r)
		return
	} else {
		http_helpers.PrintError(http.StatusNotFound, "Not Found", w)
		return
	}
}
