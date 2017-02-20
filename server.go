package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Incident struct {
	ID      int
	Date    string
	City    string
	Address string
	Lat     float64
	Lng     float64
	Killed  int
	Injured int
	News    string
}
type Device struct {
	Imei  string
	Token string
}
type Location struct {
	Imei string
	Lat  float64
	Lng  float64
}
type Range struct {
	Imei  string
	Range float64
}

func getData() (results []Incident) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shooting").C("accidents")
	err = c.Find(nil).All(&results)
	if err != nil {
		panic(err)
	}
	return
}

func saveDevice(device Device) (ok bool) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shooting").C("clients")

	_, err = c.Upsert(bson.M{"imei": device.Imei}, bson.M{"$set": bson.M{"token": device.Token}})

	if err != nil {
		panic(err)
	}
	ok = true
	return
}

func saveRange(radius Range) (ok bool) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shooting").C("clients")

	_, err = c.Upsert(bson.M{"imei": radius.Imei}, bson.M{"$set": bson.M{"range": radius.Range}})

	if err != nil {
		panic(err)
	}
	ok = true
	return
}

func saveLocation(location Location) (ok bool) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shooting").C("clients")

	_, err = c.Upsert(bson.M{"imei": location.Imei}, bson.M{"$set": bson.M{"lat": location.Lat, "lng": location.Lng}})

	if err != nil {
		panic(err)
	}
	ok = true
	return
}

// TODO: add Save Range
func main() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", func(w rest.ResponseWriter, req *rest.Request) {
			w.WriteJson(getData())
		}),
		rest.Post("/saveDevice", func(w rest.ResponseWriter, req *rest.Request) {
			device := Device{}
			err := req.DecodeJsonPayload(&device)
			if err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ok := saveDevice(device)
			if ok {
				w.WriteJson(bson.M{"saveDevice": "ok"})
			} else {
				w.WriteJson(bson.M{"saveDevice": "failed"})
			}

		}),
		rest.Post("/saveRange", func(w rest.ResponseWriter, req *rest.Request) {
			radius := Range{}
			err := req.DecodeJsonPayload(&radius)
			if err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ok := saveRange(radius)
			if ok {
				w.WriteJson(bson.M{"saveRange": "ok"})
			} else {
				w.WriteJson(bson.M{"saveRange": "failed"})
			}

		}),
		rest.Post("/saveLocation", func(w rest.ResponseWriter, req *rest.Request) {
			location := Location{}
			err := req.DecodeJsonPayload(&location)
			if err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ok := saveLocation(location)
			if ok {
				w.WriteJson(bson.M{"saveLocation": "ok"})
			} else {
				w.WriteJson(bson.M{"saveLocation": "failed"})
			}

		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":80", api.MakeHandler()))

}
