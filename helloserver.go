package main

import (
	"fmt"
	"bytes"
    "net/http"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
	"strings"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"context"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var client *firestore.Client

type lmfao struct {
	counter int
	user string
}

func sendMessage(message string) {
	fmt.Print(message) // dude figure out how to put the message variable in the jsonString
	url := "https://hooks.slack.com/services/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"
	var jsonString = []byte(`{"text": "uhhh"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

    if err != nil {
        panic(err)
	}
	
	defer resp.Body.Close()
}

func kick(channel string, user string) {
	fmt.Println("kick me bitch")
	url := "https://slack.com/api/groups.kick"
	var jsonString = []byte(`{"token": "xoxp-44822465026-431865943765-549665390069-a08f3d4819c5dcb11401d330cb944f23", "channel": "` + channel + `", "user": " "` + user +`"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(resp);

    if err != nil {
        panic(err)
	}

	defer resp.Body.Close()

}

func servePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "My localhost is online!!!!!")
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var messageData struct {
		Challenge string
		Event map[string]interface{}
		Type string
	}

	// Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
	err := decoder.Decode(&messageData) 
	
	if err != nil {
		fmt.Println("1")
		panic(err)
	}

	t := messageData.Type 

	if t == "url_verification" {
		challenge := messageData.Challenge
		fmt.Fprintf(w, string(challenge)) // sends the "challenge" response back to slack 
	} else { 

		fmt.Println("is it going in here????")

		eventDict := messageData.Event

		if eventDict == nil { return }

		msg := eventDict["text"].(string)
		userToken := eventDict["user"].(string)
		channel := eventDict["channel"].(string)
		
		fmt.Println(msg)
		fmt.Println(userToken)
		fmt.Println(channel)

		// ref := client.NewRef()
		// usersRef := ref.Child("users")
		
		if strings.Contains(strings.ToLower(msg), "uwu") {

			// update database

			/* 

			*/

			iter := client.Collection("users").Documents(context.Background())
			
			found := false

			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalf("Failed to iterate: %v", err)
				}

				data := doc.Data()

				if data["user"] == userToken {
					counter := data["counter"].(int64) + 1
					// _, _, err := client.Collection("users").Set(context.Background(), map[string]interface{}{
					// 	"counter": counter,
					// }) 

					// if err != nil {
					// 	log.Fatalln(err)
					// }


					if counter == 5 {
						kick(channel, userToken)
					}

					found = true
				}
				
				fmt.Println(doc.Data())
			}

			if !found {
				_, _, err := client.Collection("users").Add(context.Background(), map[string]interface{}{
					"counter": 0,
					"user": userToken,
				})
	
				if err != nil {
					log.Fatalln(err)
				}
			}


			defer client.Close()


			// send a message to the channel with the username and count 

		}
	}
}



func Authorize(w http.ResponseWriter, r *http.Request) {
	
}

func main() {

	kick("GG5DDUX5Y","UCPRFTRNH")

	// firebase stuff

	sa := option.WithCredentialsFile("./ServiceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)

	client, err = app.Firestore(context.Background())	

	if err != nil {
		log.Fatalln(err) 
	}

	// hi := lmfao{user: "TESTTESTTESTTEST", counter: 99}

	// result, err := client.Collection("users").Doc("user").Set(context.Background(), hi)

	// if err != nil {
	// 	log.Fatalln(err)
	// }


	// _, _, err = client.Collection("users").Add(context.Background(), map[string]interface{}{
    //     "counter": 0,
    //     "user": "anutha",
	// })

	// _, _, err = client.Collection("users").Add(context.Background(), map[string]interface{}{
    //     "counter": 0,
    //     "user": "one",
	// })

	// _, _, err = client.Collection("users").Add(context.Background(), map[string]interface{}{
    //     "counter": 0,
    //     "user": "u w u ",
	// })

	// iter := client.Collection("users").Documents(context.Background())

	// for {
    //     doc, err := iter.Next()
    //     if err == iterator.Done {
    //         break
    //     }
    //     if err != nil {
    //         log.Fatalf("Failed to iterate: %v", err)
	// 	}
		
	// 	fmt.Println(doc.Data())

	// 	// check to see if doc.data()["user"] matches the token of the person we want to update
	// 	// if it does, update the counter field 
	// 	// if counter > 5 
	// 		// kick them out of the channel
		
	// }

	// defer client.Close()


	router := mux.NewRouter()
	router.HandleFunc("/slack", GetMessage).Methods("POST")
	router.HandleFunc("/redirect", Authorize).Methods("GET")
	log.Fatal(http.ListenAndServe(":6969", router))
	
}
