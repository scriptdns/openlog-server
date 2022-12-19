package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	_ "github.com/lib/pq"
)

func main() {
	// Postgres string

	//// Postgres / Timescaledb output module
	const POSTGRES_CONNECTION_STRING = "user=postgres sslmode=disable host=localhost port=5432 user=postgres password=password database=openlog"
	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Can't ping database: " + err.Error())
	}
	stmt, err := db.Prepare("INSERT INTO \"logs\"(\"time\", stream, data) VALUES (current_timestamp, $1, $2)")
	if err != nil {
		panic(err)
	}

	//// Websocket / openlog protocol input module
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s, type %d", string(msg), mt)
			if json.Valid(bytes.TrimSpace(msg)) {
				// The message is valid json - put in sql
				_, err := stmt.Query(c.Params("id"), msg)
				if err != nil {
					log.Printf("Could not store the log data", msg)
					break
				}
				log.Printf("Stored a log: %s\n", msg)

			} else {
				log.Printf("Not storing log data because it's not json: %s\n", msg)
			}
		}

	}))

	log.Fatal(app.Listen(":3000"))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
}
