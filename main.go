package main

import (
	"context"
	"encoding/json"
	"fmt" 
	"go.mongodb.org/mongo-driver/v2/bson" 
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"	
	"net/http"
    "github.com/gin-gonic/gin"
)

type url struct {
    Site  string  `json:"site"`
    ShortKey string  `json:"shortKey"`
}

func main() {
	
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	coll := client.Database("UrlDB").Collection("Urls")
	
	router := gin.Default()	
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
	router.GET("/:shortKey", func(c *gin.Context) {
		
		shortKey := c.Param("shortKey")		
		var result bson.M
		err = coll.FindOne(context.TODO(), bson.D{{"shortkey", shortKey}}).
			Decode(&result)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the shortKey %s\n", shortKey)
			return
		}
		if err != nil {
			panic(err)
		}		
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", jsonData)
		c.Redirect(http.StatusFound, result["site"].(string))
		
	})
	router.POST("/", func(c *gin.Context) {
		
		var newUrl url
		if err := c.BindJSON(&newUrl); err != nil {
			return
		}	
		result, err := coll.InsertOne(context.TODO(), newUrl)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)	
		c.IndentedJSON(http.StatusCreated, newUrl)
		
	})
	
	router.Run("localhost:8080")
	
}