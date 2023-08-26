package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kkboranbay/hotel-reservation/api"
	"github.com/kkboranbay/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://kkboranbay:Mongodb12@cluster0.elxfrcu.mongodb.net/?retryWrites=true&w=majority"
const dbname = "hotel-reservation"
const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// // context.TODO() is a function from the Go programming language's context package.
	// // It returns a "TODO" context, which is meant to be used when you need a context
	// // but you're not sure about the specific type of context to use yet. In essence, it's a placeholder context.
	// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // The context.Background() function returns the root context for a program or request.
	// // It should be used as the starting point for creating new contexts.
	// ctx := context.Background()
	// // If you attempt to perform an operation on a database or collection that doesn't exist,
	// // MongoDB will automatically create the database or collection for you.
	// coll := client.Database(dbname).Collection(userColl)
	// // In MongoDB, a collection is conceptually similar to a table in traditional relational databases.
	// // In MongoDB, data is stored in collections, which are containers for JSON-like documents.
	// // Collections don't enforce a fixed schema, meaning that each document within a collection can have different fields and structures.

	// user := types.User{
	// 	FirstName: "Leo",
	// 	LastName:  "Ken",
	// }
	// res, err := coll.InsertOne(ctx, user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(res)
	// // &{ObjectID("64e8b68543946ffd44930c61")}

	// // This variable will be used to hold the decoded data from the MongoDB document.
	// var james types.User
	// // FindOne. It's used to retrieve a single document from the collection that matches the given query.
	// // In this case, the query is an empty BSON document (bson.M{}), which essentially means "find any document."
	// // In the context of MongoDB and the MongoDB Go driver, bson.M{} represents a BSON document in Go.
	// // BSON stands for "Binary JSON," and it's a binary-encoded format used to represent documents and data structures in MongoDB.
	// // M{}: M is a type within the bson package. It's used to represent BSON documents in
	// // a format that's convenient for constructing queries and data structures.
	// // bson.M{"name": "Alice"}
	// // .Decode(&james): This part of the code decodes the retrieved document into the james variable.
	// // The & symbol before james is used to pass a pointer to the james variable, allowing the Decode function to modify its contents.
	// if err := coll.FindOne(ctx, bson.M{}).Decode(&james); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(james)
	// // {64e8b68543946ffd44930c61 Leo Ken}

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	app.Listen(*listenAddr)
}
