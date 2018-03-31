/*
Fumble: a popular app that lets friends know when they cross paths. The client sends the users
long,lat to a post endpoint once a second, all day long. The user should be able to open the app
at any point and see the people they've crossed paths with.
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type InPacket struct {
	UserId    int `json:"userId"`
	Longitude int `json:"long"` //int for now, float for real
	Latitude  int `json:"lat"`
}

type Location struct {
	Longitude, Latitude int //int for now, float for real
}

//AppContext contains all global variables that are shared among packages
type AppContext struct {
	LocDb map[int]Location
	//DbHandler *sql.DB // db handle to mySql
}

func main() {
	locDb := map[int]Location{} //in memory cache
	/***
	dbh, err := sql.Open("mysql", "root@/fumbledb")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dbh.Close()
	***/
	globalVar := AppContext{
		LocDb: locDb,
		//DbHandler: dbh,
	}

	log.Fatal(http.ListenAndServe("localhost:3000", &globalVar))
}

func (ds *AppContext) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/fumble":
		if req.Method == "GET" {
			fmt.Fprint(w, "Fumble, a cross path service for friends\n") // return signature of the service
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	case "/fumble/location": // POST
		if req.Method == "POST" {
			var inData InPacket
			err := json.NewDecoder(req.Body).Decode(&inData)
			if err != nil {
				http.Error(w, "Error decoding JSON request body.",
					http.StatusInternalServerError)
				fmt.Fprintf(w, "inData %v\n", inData)
			}
			inId, inLoc := inData.UserId, Location{inData.Longitude, inData.Latitude}
			if _, ok := ds.LocDb[inId]; ok { //new userId
				ds.LocDb[inId] = inLoc
				//adding useId into mySql user table
				//adding indata with create timesatmp into mysql location table
			} else { //update in memory cache
				ds.LocDb[inId] = inLoc
				//adding indata with create timesatmp into mysql location table
			}

			fmt.Fprintf(w, "http %d\n", http.StatusOK)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	case "/fumble/friends":
		if req.Method == "GET" {
			out := []int{}
			for k := range ds.LocDb {
				out = append(out, k)
			}
			dataout, err := json.MarshalIndent(out, "", " ")
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}
			fmt.Fprintf(w, "{\n All users: %s\n}", string(dataout))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

	default: //should use a router package but this will do for simple case
		if strings.Contains(req.URL.Path, "/fumble/friends/") {
			if req.Method == "GET" {
				sid := strings.TrimPrefix(req.URL.Path, "/fumble/friends/")
				id, err := strconv.Atoi(sid)
				_, ok := ds.LocDb[id]
				if err != nil || id < 0 || ok == false {
					http.Error(w, "Id not valid", http.StatusBadRequest)
				} else {
					cur := ds.LocDb[id]
					out := []int{}
					for k, v := range ds.LocDb {
						//ToDo, not using exactly equal, using within a distance d
						if k != id && v.Longitude == cur.Longitude && v.Latitude == cur.Latitude {
							//found a friend
							out = append(out, k)
						}
					}
					dataout, err := json.MarshalIndent(out, "", " ")
					if err != nil {
						log.Fatalf("JSON marshaling failed: %s", err)
					}
					fmt.Fprintf(w, "{\n friends: %s\n}", string(dataout))
				}
			} else {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			}
		} else {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "http 404, %s invalid. Only /urlVal/validate is allowed.\n", req.URL)
		}
	}
}
