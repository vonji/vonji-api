package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vonji/vonji-api/api"
	"github.com/vonji/vonji-api/models"
)

//TODO status code + all responses should be JSON

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := api.GetContext()

	users := []models.User{}
	ctx.Db.Find(&users)
	for i, user := range users { //TODO There must be another way to do this
		ctx.Db.Model(&user).Association("tags").Find(&users[i].Tags)
	}

	json.NewEncoder(w).Encode(users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := api.GetContext()
	user := models.User{}

	id, err := parseUint(mux.Vars(r)["id"]) //TODO find shorter syntax

	if err != nil {
		http.Error(w, "Parameter ID is not an unsigned integer", http.StatusBadRequest)
		return
	}

	ctx.Db.First(&user, id)
	ctx.Db.Model(&user).Association("tags").Find(&user.Tags)

	if user.ID == 0 {
		http.Error(w, fmt.Sprintf("No user with ID %d found", id), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	ctx := api.GetContext()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.Db.Create(&user) //TODO check security
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	ctx := api.GetContext()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.Db.Save(&user) //TODO check security
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	ctx := api.GetContext()

	id, err := parseUint(mux.Vars(r)["id"]) //TODO find shorter syntax

	if err != nil {
		http.Error(w, "Parameter ID is not an unsigned integer", http.StatusBadRequest)
		return
	}

	user.ID = id

	ctx.Db.Delete(&user) //Soft delete
	//TODO return error if the id does not exist
}