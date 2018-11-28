package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {

	// ko := "U81b9a85bc551320d6b6c3c415554c027"

	http.HandleFunc("/callback", callback)

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func callback(w http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("mysql", "root:157953@tcp(127.0.0.1:3306)/wildfire")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	client := &http.Client{}
	bot, err := linebot.New("fc23264a0386314e124f99e253e58888", "TeD1kI257w/C58K4RcC0ax05l53WFnzNBNamLbJfBvDwHh5ogez0T9TJ1+unw5Jl0LjhWgFPmu8MUGMvUMcKmQYHQ1IQ8gf9Sy3EQUfx4lwoSVEu7j2uwbgHJHH5Hiwe6Y5eQCk4lgRygjWcW1z9yQdB04t89/1O/w1cDnyilFU=", linebot.WithHTTPClient(client))
	if err != nil {
		panic(err.Error())
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				if strings.Contains(message.Text, "รับทราบ ตัวแจ้งเตือนที่ ") {
					runes := []rune(message.Text)
					node := string(runes[24:])
					fmt.Println(node)
					edit, err := db.Prepare("update node set status = 2 where node_id = ?")
					if err != nil {
						log.Fatal(err)
					}
					defer edit.Close()
					_, err = edit.Exec(node)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Dai")
				}
			}
		} else if event.Type == linebot.EventTypeFollow {
			userID := event.Source.UserID
			fmt.Println("Follow Id = " + userID)
			insert, err := db.Prepare("insert into users(userID) values(?)")
			if err != nil {
				log.Fatal(err)
			}
			defer insert.Close()
			_, err = insert.Exec(userID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Insert UserID to Database")
		} else if event.Type == linebot.EventTypeUnfollow {
			userID := event.Source.UserID
			delete, err := db.Prepare("delete from users where userID = ?")
			if err != nil {
				log.Fatal(err)
			}
			defer delete.Close()
			_, err = delete.Exec(userID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Unfollow Id = " + userID)
			fmt.Println("Delete User From Database")
		} else if event.Type == linebot.EventTypeJoin {
			groupID := event.Source.GroupID
			fmt.Println("Group Id = " + groupID)
			insert, err := db.Prepare("insert into users(userID) values(?)")
			if err != nil {
				log.Fatal(err)
			}
			defer insert.Close()
			_, err = insert.Exec(groupID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Insert GroupID to Database")
		}
	}
}
