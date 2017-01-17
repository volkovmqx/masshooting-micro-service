package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
    "gopkg.in/mgo.v2"
    "fmt"
)
type Incident struct {
  Id int
  Date string
  City	 string
  Address string
  Cordinnates string
  Killed int
  Injured int
  News string
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
func main() {

    fmt.Println("Results All: ", )
    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)
    api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
        w.WriteJson(getData())
    }))
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
