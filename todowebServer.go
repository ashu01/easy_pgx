package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx"
    "strings"
)

// Database connectivity variables
var db *pgx.ConnPool
var db_err error

//Initialise connection to the database
func init() {
	db, db_err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Database: "tododb",
			User:     "todo",
			Password: "todo123",
			Port:     5432,
		},
		MaxConnections: 10,
	})

	if db_err != nil {
		fmt.Println("Can't connect to database")
	}
}



type Task struct {
	ID       string `json:"id,omitempty"`
	Content     string `json:"content,omitempty"`
	Complete bool   `json:"complete"`
}

type WSMessage struct{
 MessageData []Task `json:"messageData,omitempty"`
 MessageLabel string `json:"messageLabel"`
}
// type Array struct{
// 	Collection []Request
// }

type Response struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Process   string `json:"process,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/v5/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, errors = upgrader.Upgrade(w, r, nil)
		if errors != nil {
			fmt.Println(errors)
		}

		rows, err := db.Query(`SELECT * FROM todo`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var idd string
		var contentt string
		var comp bool
        var arr []Task
		for rows.Next() {
			err := rows.Scan(&idd, &contentt, &comp)
			if err != nil {
				log.Fatal(err)
			}
			var task Task
			task.ID = idd
			task.Content = contentt
			task.Complete = comp
			// fmt.Println(comp, task.Complete)
            arr = append(arr, task)
			// task_json, _ := json.Marshal(task)
			// fmt.Println(json_err)
			// fmt.Println(string(task_json))
			// conn.WriteMessage(1, task_json)
		}
        fmt.Println(arr)
        task_json, _ := json.Marshal(arr)
		// fmt.Println(json_err)
		// fmt.Println(string(task_json))
		conn.WriteMessage(1, task_json)
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn *websocket.Conn) {
			for {

				    _, msg, err := conn.ReadMessage()
				    if err != nil {
					    conn.Close()
					    fmt.Println("Web server connection closed : ", err)
                        return
				    }

                    // println(string(msg))
    				// fmt.Println(len(msg))

                    var messagereceived WSMessage
                    err1 := json.Unmarshal(msg, &messagereceived)
                    if err1 != nil {
				    	fmt.Println("error : ", err1)
				    }

                    // fmt.Println(messagereceived)
                    var operation string                 
                    operation = messagereceived.MessageLabel
                    operation = strings.ToLower(operation)

                    var datareceived []Task
                    datareceived = messagereceived.MessageData 
				    // fmt.Println(datareceived)
                    // fmt.Println(len(datareceived))
                    // fmt.Println(operation)
                    var CONTENT,IDD string
                    var COMPLETE bool

                    if operation == "insert" {
                        for _, eachstring := range datareceived {
                            CONTENT = eachstring.Content
                            IDD = eachstring.ID
                            COMPLETE = eachstring.Complete
                            tx, _ := db.Begin() // tx => transaction , _ => error and execute
						    defer tx.Rollback() // it will be executed after the completion of local function
						    _, err := tx.Exec(`
         							    INSERT INTO todo (id, task, complete) VALUES ($1, $2, $3)
    							    `, IDD, CONTENT, COMPLETE)
						
                            fmt.Println("Insert error : ", err)

    						commitErr := tx.Commit()
    						fmt.Println("Insert commit : ", commitErr)
	    					if commitErr != nil {
		    					tx.Rollback()
			    				var resp2 Response
				    			resp2.Timestamp = 501
					    		resp2.Process = "Insert commit error"
						    	resp_json, _ := json.Marshal(resp2)
							    // fmt.Println(json_err)
    							// fmt.Println(string(resp_json))
	    						conn.WriteMessage(1, resp_json)
		    					// conn.WriteJSON(Response{
			    				// 	Timestamp: 123445,
				    			// })
					    		break
						    }
                        }                  
                    } else if operation == "delete" {
                        for _, eachstring := range datareceived {
                            CONTENT = eachstring.Content
                            IDD = eachstring.ID
                            COMPLETE = eachstring.Complete
                            tx, _ := db.Begin() // tx => transaction , _ => error and execute
    						defer tx.Rollback() // it will be executed after the completion of local function
                            _, err := tx.Exec(`
             	    					DELETE FROM todo WHERE  id = $1
    		    					`, IDD)
					    	fmt.Println("Delete : ", err)

						    commitErr := tx.Commit()
    						fmt.Println("Delete commit :  ", commitErr)
	    					if commitErr != nil {
		    					tx.Rollback()
			    				var resp1 Response
				    			resp1.Timestamp = 501
					    		resp1.Process = "Delete Commit Error"
						    	resp_json, _ := json.Marshal(resp1)
							    // fmt.Println(json_err)
    							// fmt.Println(string(resp_json))
	    						conn.WriteMessage(1, resp_json)
		    					// conn.WriteJSON(Response{
			    				// 	Timestamp: 123445,
				    			// })
					    		break
						    }
                        }                
                    } else if operation == "update" {
                        for _, eachstring := range datareceived {
                            CONTENT = eachstring.Content
                            IDD = eachstring.ID
                            COMPLETE = eachstring.Complete
                            
                            tx, _ := db.Begin() // tx => transaction , _ => error and execute
    						defer tx.Rollback() // it will be executed after the completion of local function
	    					_, err := tx.Exec(`
             	    					UPDATE todo SET complete = $1 WHERE id = $2
    		    					`, COMPLETE, IDD)
					    	fmt.Println("Update : ", err)

						    commitErr := tx.Commit()
    						fmt.Println("Update commit :  ", commitErr)
	    					if commitErr != nil {
		    					tx.Rollback()
			    				var respu Response
				    			respu.Timestamp = 501
					    		respu.Process = "Update Commit Error"
						    	resp_json, _ := json.Marshal(respu)
							    // fmt.Println(json_err)
    							// fmt.Println(string(resp_json))
	    						conn.WriteMessage(1, resp_json)
		    					// conn.WriteJSON(Response{
			    				// 	Timestamp: 123445,
				    			// })
					    		break
						    }
                        }                
                  } else {
                        fmt.Println("No such Process : ", operation)
                 }
			    	var response Response
					response.Timestamp = 501
					response.Process = operation + " Success"
					response_json, _ := json.Marshal(response)
					// fmt.Println(json_err)
					fmt.Println(string(response_json))
					conn.WriteMessage(1, response_json)
			}
		}(conn)
	})

	fmt.Println("Live on :3000")
	http.ListenAndServe(":3000", nil)
}
