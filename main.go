package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tabrizihamid84/mongo-golang/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserControllers(getSession())
	r.GET("/users", uc.GetUsers)
	r.GET("/users/:id", uc.GetUser)
	r.POST("/users", uc.CreateUser)
	r.DELETE("/users/:id", uc.DeleteUser)
	fmt.Println("server running at port 9000")
	http.ListenAndServe("localhost:9000", r)

}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	return s
}
