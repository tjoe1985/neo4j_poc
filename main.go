package main

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func main() {
	uri := "neo4j://localhost:7687"
	pass := "todospassword"
	usr := "neo4j"
	//db := "todos"

	result, err := helloWorld(context.Background(), uri, usr, pass, "users")
	if err != nil {
		log.Println("Error from hello World: ", err)
	}

	log.Println("results: ", result)

}

func helloWorld(ctx context.Context, uri, username, password string, db string) ([]User, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: db})
	defer session.Close(ctx)

	//make empty array of users
	usrs := []User{}

	_, err = session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			"MATCH (n) RETURN n.email as email, n.firstName as name, n.lastName as lastname;",
			//map[string]any{}
			nil)
		if err != nil {
			log.Println("Error running trasaction: ", err)
			return nil, err
		}

		for result.Next(ctx) {

			record := result.Record()

			Email, err := record.Get("email")
			if !err {
				log.Println("Error getting email: ", err)
			}
			FirstName, err := record.Get("name")
			if !err {
				log.Println("Error getting name: ", err)
			}
			LastName, err := record.Get("lastname")
			if !err {
				log.Println("Error getting lastname: ", err)
			}
			s := User{
				Email:     Email.(string),
				FirstName: FirstName.(string),
				LastName:  LastName.(string),
			}

			usrs = append(usrs, s)

		}

		return usrs, result.Err()
	})
	if err != nil {
		log.Println("Error executing write: ", err)
		return usrs, err
	}

	return usrs, nil
}
