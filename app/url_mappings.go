package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
	"github.com/piyushm-dt/zusers/db"
	"github.com/piyushm-dt/zusers/domain/employee"
)

var collection = db.ConnectDB()

func mapUrls() {

	router.HandleFunc("/api/employees", createEmployee).Methods("POST")
	router.HandleFunc("/api/employees", getEmployeess).Methods("GET")
	router.HandleFunc("/api/employees/{id}", getEmployee).Methods("GET")
	router.HandleFunc("/api/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/api/employees/{id}", deleteEmployee).Methods("DELETE")
}


func getEmployeess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	
	var emp []employee.Employee

	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println("Error: ", err)
	
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		
		var emp_ employee.Employee
	
		err := cur.Decode(&emp_)
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		emp = append(emp, emp_)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(emp)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var emp employee.Employee
	
	var params = mux.Vars(r)

	
	id, _ := primitive.ObjectIDFromHex(params["id"])

	
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&emp)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	json.NewEncoder(w).Encode(emp)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emp employee.Employee

	_ = json.NewDecoder(r.Body).Decode(&emp)

	// insert our emp model.
	result, err := collection.InsertOne(context.TODO(), emp)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var emp employee.Employee

	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&emp)

	// prepare update model.
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "firstname", Value: emp.FirstName},
			{Key: "lastname", Value: emp.LastName},
			{Key: "email", Value: emp.Email},
			{Key: "designation", Value: emp.Designation},
			{Key: "skills", Value: emp.Skills},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&emp)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	emp.ID = id

	json.NewEncoder(w).Encode(emp)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}