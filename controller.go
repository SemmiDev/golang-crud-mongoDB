package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorilla/mux"
	"net/http"
)

type student struct {
	Nisn string 	`json:"nisn"`
	Name string 	`json:"name"`
	Age int	 		`json:"age"`
	Address string 	`json:"address"`
}

// context
var Context = context.TODO()

// get collection students from database belajar_golang
var userCollection = database().Database("belajar_golang").Collection("students")

// create student profile
func createProfile(writer http.ResponseWriter, request *http.Request) {
	// content type
	writer.Header().Set("Content-Type", "application/json")

	var person student
	// kita store variable person of type student
	err := json.NewDecoder(request.Body).Decode(&person)
	if err != nil {
		fmt.Print(err)
	}

	insertResult, err := userCollection.InsertOne(Context, person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("memasukkan 1 dokumen -> ", insertResult)
	//mengembalikan id dari dokumen yg di bikin
	json.NewEncoder(writer).Encode(insertResult.InsertedID)
}

// get profile of a particular student by name
func getStudentProfile(writer http.ResponseWriter, request *http.Request) {
	// content type
	writer.Header().Set("Content-Type", "application/json")

	var body student
	// kita store variable person of type student
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		fmt.Print(err)
	}

	var result primitive.M
	err = userCollection.FindOne(Context, bson.D{{"name", body.Name}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	// returns a Map containing document
	json.NewEncoder(writer).Encode(result)
}

// update profile dari student
func updateProfile(writer http.ResponseWriter, request *http.Request) {
	// content type
	writer.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Name string 	`json:"name"` // match
		Address string 	`json:"address"` // modified
	}

	var body updateBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		fmt.Print(err)
	}

	filter := bson.D{{"name", body.Name}}
	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"address", body.Address}}}}
	updateResult := userCollection.FindOneAndUpdate(Context, filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(writer).Encode(result)
}

func deleteProfile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	//get Parameter value as string
	params := mux.Vars(request)["id"]

	// convert params to mongodb Hex ID
	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		fmt.Printf(err.Error())
	}

	opts := options.Delete().SetCollation(&options.Collation{})
	res, err := userCollection.DeleteOne(Context, bson.D{{"_id", _id}}, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	json.NewEncoder(writer).Encode(res.DeletedCount) // return number of documents deleted

}

func getAllStudents(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var results []primitive.M
	// mengembalikan sebuah *mongo.Cursor
	cur, err := userCollection.Find(Context, bson.D{{}})
	if err != nil {
		fmt.Println(err)
	}

	// kita loop
	for cur.Next(Context) {
		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		// appending document pointed by Next()
		results = append(results, elem)
	}
	// close the cursor once stream of documents has exhausted
	cur.Close(Context)
	json.NewEncoder(writer).Encode(results)
}

