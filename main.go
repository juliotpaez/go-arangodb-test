package main

import (
	"fmt"
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserRequest struct {
	UserId   string `json:"userId" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type DbLog struct {
	UserId    string
	Timestamp int64
}

func main() {
	// Start DB connection.
	connection, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		fmt.Println("Cannot connect to DB: " + err.Error())
		return
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     connection,
		Authentication: driver.BasicAuthentication("root", "test"),
	})
	if err != nil {
		fmt.Println("Cannot connect to DB: " + err.Error())
		return
	}

	db, err := openDB(client)
	if err != nil {
		return
	}

	collection, err := openCollection(db)
	if err != nil {
		return
	}

	// Initiate server.
	router := gin.Default()
	router.POST("/echo", func(c *gin.Context) {
		// Validate input
		var input UserRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Add log to DB.
		log := DbLog{
			UserId:    input.UserId,
			Timestamp: time.Now().Unix(),
		}
		_, err = collection.CreateDocument(nil, log)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello %s!", input.Username),
		})
	})
	router.Run()
}

func openDB(client driver.Client) (driver.Database, error) {
	dbExists, err := client.DatabaseExists(nil, "MicroService")
	if err != nil {
		fmt.Println("Cannot open database \"MicroService\": " + err.Error())
		return nil, err
	}

	if dbExists {
		db, err := client.Database(nil, "MicroService")
		if err != nil {
			fmt.Println("Cannot open database \"MicroService\": " + err.Error())
			return nil, err
		}

		return db, nil
	} else {
		options := &driver.CreateDatabaseOptions{}
		db, err := client.CreateDatabase(nil, "MicroService", options)
		if err != nil {
			fmt.Println("Cannot create database \"MicroService\": " + err.Error())
			return nil, err
		}

		return db, err
	}
}

func openCollection(database driver.Database) (driver.Collection, error) {
	collectionExists, err := database.CollectionExists(nil, "Logs")
	if err != nil {
		fmt.Println("Cannot open collection \"Logs\": " + err.Error())
		return nil, err
	}

	if collectionExists {
		db, err := database.Collection(nil, "Logs")
		if err != nil {
			fmt.Println("Cannot open collection \"Logs\": " + err.Error())
			return nil, err
		}

		return db, nil
	} else {
		options := &driver.CreateCollectionOptions{}
		db, err := database.CreateCollection(nil, "Logs", options)
		if err != nil {
			fmt.Println("Cannot create collection \"Logs\": " + err.Error())
			return nil, err
		}

		return db, err
	}
}
