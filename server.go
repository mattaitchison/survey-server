package main

import (
	"encoding/json"

	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stevemurr/survey-server/filedb"
)

const host = "/Users/murr/data/db.json"

var db = filedb.DB{Host: host}

func RecordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		db.Lock() // Have to lock when using db.All() to prevent data race.
		defer db.Unlock()
		records := db.All()
		results := make([]map[string]string, len(records))
		idx := 0
		for _, value := range records {
			results[idx] = value
			idx++
		}
		b, _ := json.Marshal(results)
		w.Write(b)
	case "POST":
		name := r.FormValue("name")
		age := r.FormValue("age")
		occupation := r.FormValue("occupation")
		company := r.FormValue("company")
		location := r.FormValue("location")
		discoveryCategory := r.FormValue("discoveryCategory")
		discoveryNotes := r.FormValue("discoveryNotes")
		useCategory := r.FormValue("useCategory")
		useNotes := r.FormValue("useNotes")
		email := r.FormValue("email")
		record := make(filedb.Record)
		record["name"] = name
		record["age"] = age
		record["occupation"] = occupation
		record["company"] = company
		record["location"] = location
		record["discoveryCategory"] = discoveryCategory
		record["discoveryNotes"] = discoveryNotes
		record["useCategory"] = useCategory
		record["useNotes"] = useNotes
		record["email"] = email
		r := db.InsertRecord(record)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", r)

	}
}

func RecordHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	switch r.Method {
	case "GET":
		record := db.FindID(id)
		b, _ := json.Marshal(record)
		fmt.Fprintf(w, string(b))
	case "DELETE":
		db.RemoveID(id)
		fmt.Fprintf(w, "Deleted %s", id)
	case "PUT":
		name := r.FormValue("name")
		age := r.FormValue("age")
		occupation := r.FormValue("occupation")
		company := r.FormValue("company")
		location := r.FormValue("location")
		discoveryCategory := r.FormValue("discoveryCategory")
		discoveryNotes := r.FormValue("discoveryNotes")
		useCategory := r.FormValue("useCategory")
		useNotes := r.FormValue("useNotes")
		email := r.FormValue("email")
		record := make(map[string]string)
		record["name"] = name
		record["age"] = age
		record["occupation"] = occupation
		record["company"] = company
		record["location"] = location
		record["discoveryCategory"] = discoveryCategory
		record["discoveryNotes"] = discoveryNotes
		record["useCategory"] = useCategory
		record["useNotes"] = useNotes
		record["email"] = email
		record["id"] = id
		db.UpdateID(id, record)
	}
}

var index = template.Must(template.ParseFiles("templates/index.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, "")
}

func AddRecordIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/records/index.html")
	t.Execute(w, "")
}

func main() {
	fmt.Println(banner)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/add", AddRecordIndex)
	router.HandleFunc("/records", RecordsHandler)
	router.HandleFunc("/records/{id}", RecordHandler)

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		log.Println(route.GetPathTemplate())
		return nil
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}

const banner = `

   _____                                  _____
  / ____|                                / ____|
 | (___  _   _ _ ____   _____ _ __ _   _| (___   ___ _ ____   _____ _ __
  \___ \| | | | '__\ \ / / _ \ '__| | | |\___ \ / _ \ '__\ \ / / _ \ '__|
  ____) | |_| | |   \ V /  __/ |  | |_| |____) |  __/ |   \ V /  __/ |
 |_____/ \__,_|_|    \_/ \___|_|   \__, |_____/ \___|_|    \_/ \___|_|
                                    __/ |
                                   |___/
`
