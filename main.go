package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database Structure
type Documents struct {
	ID      string    `json:"id,omitempty" bson:"_id,omitempty"`
	Need    int       `json:"need,omitempty" bson:"need,omitempty"`
	Subs    int       `json:"subs,omitempty" bson:"subs,omitempty"`
	Goal    int       `json:"goal,omitempty" bson:"goal,omitempty"`
	Created time.Time `json:"created,omitempty" bson:"created,omitempty"`
}

func main() {
	for {

		if !IsOnline() {
			log.Println("Connected to internet. Welcome to the Party")
			break
		}
	}

	connect()
	handleRequest()
}

func IsOnline() bool {
	_, err := http.Get("https://icanhazip.com/")
	if err == nil {
		return false
	}
	return true
}

var client *mongo.Client

// Connecting with the database (MongoDB)
func connect() {

	clientOptions := options.Client().ApplyURI("mongo_url")
	client, _ = mongo.NewClient(clientOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected to MondoDB Server")
	}

}

// function for handling the request from the client.
func handleRequest() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/api", homePage)

	http.HandleFunc("/api/needs", handleNeeds)
	http.HandleFunc("/api/sub", handleSubs)
	http.HandleFunc("/api/goal", handleGoal)

	err := http.ListenAndServe(":5284", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Baby welcome to the party.")
	log.Println("[RIN-DB]:: Home Page")
}

func handleNeeds(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var records []Documents
		collection := client.Database("databaseData").Collection("database")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		cursor, err := collection.Find(ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var record Documents
			cursor.Decode(&record)

			records = append(records, record)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: GET handleNeeds")

		filter := bson.D{{Key: "need", Value: records[0].Need}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "need", Value: records[0].Need + 1},
			}},
		}
		collection.UpdateOne(context.TODO(), filter, update)

		var _records []Documents
		_collection := client.Database("databaseData").Collection("database")
		_ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		_cursor, err := _collection.Find(_ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer _cursor.Close(_ctx)

		for _cursor.Next(_ctx) {
			var record Documents
			_cursor.Decode(&record)

			_records = append(_records, record)
		}
		if err = _cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: GET handleNeeds")
		json.NewEncoder(response).Encode(_records)

	} else {
		var _records []Documents
		_collection := client.Database("databaseData").Collection("database")
		_ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		_cursor, err := _collection.Find(_ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer _cursor.Close(_ctx)

		for _cursor.Next(_ctx) {
			var record Documents
			_cursor.Decode(&record)

			_records = append(_records, record)
		}
		if err = _cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: GET handleNeeds")
		json.NewEncoder(response).Encode(_records)

	}
}

func handleGoal(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var records []Documents
		collection := client.Database("databaseData").Collection("database")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		cursor, err := collection.Find(ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var record Documents
			cursor.Decode(&record)

			records = append(records, record)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: GET records")
		json.NewEncoder(response).Encode(records)

	} else {

		q := request.FormValue("goal")
		i, _ := strconv.Atoi(q)

		var records []Documents
		collection := client.Database("databaseData").Collection("database")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		cursor, err := collection.Find(ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var record Documents
			cursor.Decode(&record)
			records = append(records, record)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		filter := bson.D{{Key: "goal", Value: records[1].Goal}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "goal", Value: i},
			}},
		}

		collection.UpdateOne(context.TODO(), filter, update)
		log.Println("[RIN-DB]:: POST records")

		var _records []Documents
		_collection := client.Database("databaseData").Collection("database")
		_ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		_cursor, err := _collection.Find(_ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer _cursor.Close(_ctx)

		for _cursor.Next(_ctx) {
			var record Documents
			_cursor.Decode(&record)

			_records = append(_records, record)
		}
		if err = _cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: GET _records")
		json.NewEncoder(response).Encode(_records)

	}
}

func handleSubs(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var records []Documents
		collection := client.Database("databaseData").Collection("database")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		cursor, err := collection.Find(ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var record Documents
			cursor.Decode(&record)

			records = append(records, record)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		log.Println("[RIN-DB]:: GET records")
		json.NewEncoder(response).Encode(records)

	} else {

		q := request.FormValue("sub")
		i, _ := strconv.Atoi(q)
		var records []Documents
		collection := client.Database("databaseData").Collection("database")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		cursor, err := collection.Find(ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var record Documents
			cursor.Decode(&record)
			records = append(records, record)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		filter := bson.D{{Key: "subs", Value: records[1].Subs}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "subs", Value: records[1].Subs + i},
			}},
		}

		collection.UpdateOne(context.TODO(), filter, update)
		log.Println("[RIN-DB]:: POST records")

		var _records []Documents
		_collection := client.Database("databaseData").Collection("database")
		_ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		_cursor, err := _collection.Find(_ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer _cursor.Close(_ctx)

		for _cursor.Next(_ctx) {
			var record Documents
			_cursor.Decode(&record)

			_records = append(_records, record)
		}
		if err = _cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: GET _records")
		json.NewEncoder(response).Encode(_records)

	}
}

func returnRecords(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var records []Documents
		collection := client.Database("databaseData").Collection("database")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		cursor, err := collection.Find(ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var record Documents
			cursor.Decode(&record)
			records = append(records, record)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: GET records")
		json.NewEncoder(response).Encode(records)

	} else {
		var records []Documents
		collection := client.Database("databaseData").Collection("database")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		cursor, err := collection.Find(ctx, bson.D{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var record Documents
			cursor.Decode(&record)
			records = append(records, record)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		log.Println("[RIN-DB]:: POST records")
		json.NewEncoder(response).Encode(records)
		// insertArticle(newArticle)
	}
}

// For querying articles on id
func returnSingleArticle(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	var id string = request.URL.Path
	id = id[10:]
	var record Documents
	collection := client.Database("databaseData").Collection("database")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, Documents{ID: id}).Decode(&record)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	fmt.Println("Returned record ID NO : ", record.ID)
	json.NewEncoder(response).Encode(record)
}

// For query the database using the search query q=
func returnSearchResult(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Searching.....")
	q := request.FormValue("q")
	collection := client.Database("databaseData").Collection("database")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// creating a new index for search query
	index := mongo.IndexModel{Keys: bson.M{"title": "text", "content": "text", "subtitle": "text"}}
	if _, err := collection.Indexes().CreateOne(ctx, index); err != nil {
		log.Println("Could not create index:", err)
	}
	cursor, err := collection.Find(ctx, bson.M{"$text": bson.M{"$search": q}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	var records []Documents
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var record Documents
		cursor.Decode(&record)
		records = append(records, record)
	}
	if err = cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	fmt.Println("Endpoint Hit: returnAllrecords")
	json.NewEncoder(response).Encode(records)
}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
