package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	//"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/piyushm-dt/zusers/db"
	"github.com/piyushm-dt/zusers/domain/auth"
	"github.com/piyushm-dt/zusers/domain/employee"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	pong="Server is Alive!"
)
var (
	collection = db.ConnectDB()
	secretkey string = "monGooSe"
)

func mapUrls() {

	router.HandleFunc("/", ping).Methods("GET")
	router.HandleFunc("/api/signin", SignIn).Methods("POST")

	//router.HandleFunc("/api/employees", createEmployee).Methods("POST")
	router.HandleFunc("/api/employees", IsAuthorized(CreateRoleIndex)).Methods("POST")

	//router.HandleFunc("/api/employees", getEmployeess).Methods("GET")
	router.HandleFunc("/api/employees", IsAuthorized(GetRoleIndex)).Methods("GET")
	//router.HandleFunc("/api/employees/pagination", getEmployeessPagination).Methods("GET")
	router.HandleFunc("/api/employees/pagination", IsAuthorized(PaginationGetRoleIndex)).Methods("GET")
	//router.HandleFunc("/api/employees/{id}", getEmployee).Methods("GET")
	router.HandleFunc("/api/employees/{id}", IsAuthorized(IdGetRoleIndex)).Methods("GET")
	//router.HandleFunc("/api/employees/name/{firstname}", getEmployeeByName).Methods("GET")
	router.HandleFunc("/api/employees/name/{firstname}", IsAuthorized(NameGetRoleIndex)).Methods("GET")

	//router.HandleFunc("/api/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/api/employees/{id}", IsAuthorized(UpdateRoleIndex)).Methods("PUT")
	
	//router.HandleFunc("/api/employees/{id}", deleteEmployee).Methods("DELETE")
	router.HandleFunc("/api/employees/{id}", IsAuthorized(DeleteRoleIndex)).Methods("DELETE")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pong))
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
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

func getEmployeesPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var emp []employee.Employee
	
	page_size, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page_size < 1 {
		page_size = 1
	}
	page_limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || page_limit < 1 {
		page_limit = 1
	}
	
	skips := page_limit * (page_size - 1)

	opt := options.FindOptions{}
	opt.SetSkip(int64(skips))
	opt.SetLimit(int64(page_limit))

	cur, err := collection.Find(context.TODO(), bson.M{}, &opt)
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
		emp = append(emp, emp_)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(emp)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
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

func getEmployeeByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emp []employee.Employee

	cur, err := collection.Find(context.TODO(), bson.M{"firstname": mux.Vars(r)["firstname"]})
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

func createEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emp employee.Employee
	_ = json.NewDecoder(r.Body).Decode(&emp)

	emp.Password,_ = GeneratehashPassword(emp.Password)
	
	result, err := collection.InsertOne(context.TODO(), emp)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emp employee.Employee
	var params = mux.Vars(r) // using mux get variables from request URL like {id}
	
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&emp)

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

func SignIn(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var authdetails auth.Authentication
	errZ := json.NewDecoder(r.Body).Decode(&authdetails)
	if errZ != nil {
		err := "Invalid header credentials"
		json.NewEncoder(w).Encode(err)
		return
	}

	var emp employee.Employee 
	err := collection.FindOne(context.TODO(), bson.M{"email": authdetails.Email}).Decode(&emp)
	if err != nil {
		json.NewEncoder(w).Encode("No such email exists!")
		return
	}
	if emp.Email == "" {
		err := "Invalid email credentials"
		json.NewEncoder(w).Encode(err)
		return
	}

	check := CheckPasswordHash(authdetails.Password, emp.Password)
	if check {
		json.NewEncoder(w).Encode("Invalid password credentials")
		return
	}

	validToken, err := GenerateJWT(emp.Email, emp.Designation)
	if err != nil {
		json.NewEncoder(w).Encode("Invalid jwt token credentials")
		return
	}
	var token auth.Token
	token.Email = emp.Email
	token.Designation = emp.Designation
	token.TokenString = validToken

	json.NewEncoder(w).Encode(token)
}

func GenerateJWT(email, designation string) (string, error) {
	var mySignedKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["designation"] = designation
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := token.SignedString(mySignedKey)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func CheckPasswordHash(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc (func (w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			json.NewEncoder(w).Encode("request denied due to no access token")
			return
		}

		var mySignedKey = []byte(secretkey)
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token)(interface{}, error){
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Error")
			}
			return mySignedKey, nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(errors.New("access token expired"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["designation"] == "admin" {
				r.Header.Set("Designation", "admin")
				handler.ServeHTTP(w, r)
			} else if claims["designation"] == "hr" {
				r.Header.Set("Designation", "hr")
				handler.ServeHTTP(w, r)
			} else if claims["designation"] == "hradmin" {
				r.Header.Set("Designation", "hradmin")
				handler.ServeHTTP(w, r)
			} else {
				json.NewEncoder(w).Encode("you are not allowed to do that")
			}
		} else {
			erry := errors.New("invalid access")
			json.NewEncoder(w).Encode(erry)
		}
	})
}

func DeleteRoleIndex(w http.ResponseWriter, r *http.Request){

	if r.Header.Get("Designation") == "admin" {
		deleteEmployee(w, r)
	} else {
		w.Write([]byte("Not Authorized"))
	}
}

func UpdateRoleIndex(w http.ResponseWriter, r *http.Request){

	if r.Header.Get("Designation") == "admin" {
		updateEmployee(w, r)
	} else if r.Header.Get("Designation") == "hradmin" {
		updateEmployee(w, r)
	} else {
		w.Write([]byte("Not Authorized"))
	}
		
}

func CreateRoleIndex(w http.ResponseWriter, r *http.Request){
	if r.Header.Get("Designation") == "admin" {
		createEmployee(w, r)
	} else if r.Header.Get("Designation") == "hradmin" {
		createEmployee(w, r)
	} else {
		w.Write([]byte("Not Authorized"))
	}
}

func GetRoleIndex(w http.ResponseWriter, r *http.Request){
	if r.Header.Get("Designation") == "admin" {
		getEmployees(w, r)
	} else if r.Header.Get("Designation") == "hradmin" {
		getEmployees(w, r)
	} else if r.Header.Get("Designation") == "hr" {
		getEmployees(w, r)
	} else {
		w.Write([]byte("Not Authorized"))
	}
}

func IdGetRoleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Designation") == "admin" {
		getEmployee(w, r)
	} else if r.Header.Get("Designation") == "hradmin" {
		getEmployee(w, r)
	} else if r.Header.Get("Designation") == "hr" {
		getEmployee(w, r)
	} else {
		w.Write([]byte("Not Authorized"))
	}
}

func NameGetRoleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Designation") == "admin" {
		getEmployeeByName(w, r)
	} else if r.Header.Get("Designation") == "hradmin" {
		getEmployeeByName(w, r)
	} else if r.Header.Get("Designation") == "hr" {
		getEmployeeByName(w, r)
	} else {
		w.Write([]byte("Not Authorized"))
	}
}

func PaginationGetRoleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Designation") == "admin" {
		getEmployeesPagination(w, r)
	} else if r.Header.Get("Designation") == "hradmin" {
		getEmployeesPagination(w, r)
	} else if r.Header.Get("Designation") == "hr" {
		getEmployeesPagination(w, r)
	} else {
		w.Write([]byte("Not Authorized"))
	}
}
