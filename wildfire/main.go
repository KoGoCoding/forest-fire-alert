package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	r := gin.Default()
	r.POST("/savedata", savedata)
	r.POST("/alert", alert)
	// http.HandleFunc("/callback", callback)
	// if err := http.ListenAndServe(":8081", nil); err != nil {
	// 	log.Fatal(err)
	// }

	r.Run()
}

func savedata(c *gin.Context) {
	co2 := c.PostForm("co2")
	temp := c.PostForm("temp")
	// humi := c.PostForm("humi")
	nodeID := c.PostForm("node")

	db, err := sql.Open("mysql", "root:157953@tcp(127.0.0.1:3306)/wildfire")
	if err != nil {
		panic(err.Error)
	}
	defer db.Close()

	insert, err := db.Prepare("insert into data(co2, temperature, node) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insert.Close()
	_, err = insert.Exec(co2, temp, nodeID)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"CO2":         co2,
		"Temperature": temp,
	})

}

func alert(c *gin.Context) {
	db, err := sql.Open("mysql", "root:157953@tcp(127.0.0.1:3306)/wildfire")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var userID string
	// err = db.QueryRow("select * from node where node_id = ?", 1).Scan(&node1.id, &node1.lat, &node1.long)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var userIDs []string

	result, err := db.Query("select userId from users")
	if err != nil {
		log.Fatal(err)
	}

	for result.Next() {
		err = result.Scan(&userID)
		if err != nil {
			log.Fatal(err)
		}
		userIDs = append(userIDs, userID)
	}

	client := &http.Client{}
	bot, err := linebot.New("fc23264a0386314e124f99e253e58888", "TeD1kI257w/C58K4RcC0ax05l53WFnzNBNamLbJfBvDwHh5ogez0T9TJ1+unw5Jl0LjhWgFPmu8MUGMvUMcKmQYHQ1IQ8gf9Sy3EQUfx4lwoSVEu7j2uwbgHJHH5Hiwe6Y5eQCk4lgRygjWcW1z9yQdB04t89/1O/w1cDnyilFU=", linebot.WithHTTPClient(client))
	if err != nil {
		panic(err.Error())
	}

	type Node struct {
		latitude  float64
		longitude float64
	}

	type Data struct {
		co2         string
		temperature string
		humidity    string
	}

	var node Node
	id := "1"
	err = db.QueryRow("select latitude, longitude from node where node_id = ?", id).Scan(&node.latitude, &node.longitude)
	if err != nil {
		log.Fatal(err)
	}

	var data Data
	err = db.QueryRow("select co2, temperature, humidity from data where data.node = ? order by add_time desc limit 1", "1").Scan(&data.co2, &data.temperature, &data.humidity)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()
	// i, err := strconv.Atoi(data.co2)
	// var co2Status string
	// if i < 2300 {
	// 	co2Status = "น้อย"
	// } else if i < 2900 {
	// 	co2Status = "พบเล็กน้อย"
	// } else if i > 2900 {
	// 	co2Status = "มีควัน"
	// }

	// bot.PushMessage("U81b9a85bc551320d6b6c3c415554c027", linebot.NewTemplateMessage("แจ้งเตือนไฟป่า", linebot.NewButtonsTemplate("", "แจ้งเตือนไฟป่า\r\nตัวแจ้งเตือนที่ "+id, "CO2 = "+data.co2+" ppm\r\nTemperture = "+data.temperature+" *C", linebot.NewMessageAction("รับทราบ ตัวแจ้งเตือนที่ "+id, "รับทราบ ตัวแจ้งเตือนที่ "+id)))).Do()
	// bot.PushMessage("U81b9a85bc551320d6b6c3c415554c027", linebot.NewLocationMessage("แจ้งเตือนไฟป่า", "ตัวแจ้งเตือนที่ "+id, node.latitude, node.longitude)).Do()
	bot.Multicast(userIDs, linebot.NewTemplateMessage("แจ้งเตือนไฟป่า", linebot.NewButtonsTemplate("", "แจ้งเตือนไฟป่า\r\nตัวแจ้งเตือนที่ "+id, "พื้นที่เสี่ยงภัยเกิดไฟป่า\r\n", linebot.NewMessageAction("รับทราบ ตัวแจ้งเตือนที่ "+id, "รับทราบ ตัวแจ้งเตือนที่ "+id)))).Do()
	bot.Multicast(userIDs, linebot.NewLocationMessage("แจ้งเตือนไฟป่า", "ตัวแจ้งเตือนที่ "+id, node.latitude, node.longitude)).Do()
	// CO2 = 150 ppm\r\nTemperture = 40 *C\r\nHumidity = 30%

}

// func callback(w http.ResponseWriter, req *http.Request) {
// 	db, err := sql.Open("mysql", "root:157953@tcp(127.0.0.1:3306)/wildfire")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	client := &http.Client{}
// 	bot, err := linebot.New("fc23264a0386314e124f99e253e58888", "TeD1kI257w/C58K4RcC0ax05l53WFnzNBNamLbJfBvDwHh5ogez0T9TJ1+unw5Jl0LjhWgFPmu8MUGMvUMcKmQYHQ1IQ8gf9Sy3EQUfx4lwoSVEu7j2uwbgHJHH5Hiwe6Y5eQCk4lgRygjWcW1z9yQdB04t89/1O/w1cDnyilFU=", linebot.WithHTTPClient(client))
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	events, err := bot.ParseRequest(req)
// 	if err != nil {
// 		if err == linebot.ErrInvalidSignature {
// 			w.WriteHeader(400)
// 		} else {
// 			w.WriteHeader(500)
// 		}
// 		return
// 	}
// 	for _, event := range events {
// 		if event.Type == linebot.EventTypeMessage {
// 			switch message := event.Message.(type) {
// 			case *linebot.TextMessage:

// 				if strings.Contains(message.Text, "รับทราบ ตัวแจ้งเตือนที่ ") {
// 					runes := []rune(message.Text)
// 					node := string(runes[24:])
// 					fmt.Println(node)
// 					edit, err := db.Prepare("update node set status = 2 where node_id = ?")
// 					if err != nil {
// 						log.Fatal(err)
// 					}
// 					defer edit.Close()
// 					_, err = edit.Exec(node)
// 					if err != nil {
// 						log.Fatal(err)
// 					}
// 					fmt.Println("Dai")
// 				}
// 			}
// 		} else if event.Type == linebot.EventTypeFollow {
// 			userID := event.Source.UserID
// 			fmt.Println("Follow Id = " + userID)
// 			insert, err := db.Prepare("insert into users(userID) values(?)")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer insert.Close()
// 			_, err = insert.Exec(userID)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			fmt.Println("Insert UserID to Database")
// 		} else if event.Type == linebot.EventTypeUnfollow {
// 			userID := event.Source.UserID
// 			delete, err := db.Prepare("delete from users where userID = ?")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer delete.Close()
// 			_, err = delete.Exec(userID)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			fmt.Println("Unfollow Id = " + userID)
// 			fmt.Println("Delete User From Database")
// 		} else if event.Type == linebot.EventTypeJoin {
// 			groupID := event.Source.GroupID
// 			fmt.Println("Group Id = " + groupID)
// 			insert, err := db.Prepare("insert into users(userID) values(?)")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer insert.Close()
// 			_, err = insert.Exec(groupID)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			fmt.Println("Insert GroupID to Database")
// 		}
// 	}
// }
