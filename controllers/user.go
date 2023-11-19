package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tabrizihamid84/mongo-golang/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserControllers(s *mgo.Session) *UserController {
	return &UserController{
		session: s,
	}
}

func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var users []models.User

	if err := uc.session.DB("mongo-golang").C("users").Find(bson.M{}).All(&users); err != nil {
		w.WriteHeader(404)
		return
	}

	uj, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	u := models.User{}

	if err := uc.session.DB("mongo-golang").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.Id = bson.NewObjectId()

	uc.session.DB("mongo-golang").C("users").Insert(u)

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)

}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	var user models.User

	if err := uc.session.DB("mongo-golang").C("users").FindId(oid).One(&user); err != nil {
		w.WriteHeader(404)
		return
	}

	if err := uc.session.DB("mongo-golang").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user %s\n", bson.M{id: oid})
}
