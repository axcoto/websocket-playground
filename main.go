package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

type Room struct {
	members []chan string
}

func (r *Room) Join(ch chan string) {
	r.members = append(r.members, ch)
}

func (r *Room) Send(msg []byte) {
	for _, ch := range r.members {
		ch <- string(msg)
	}
}

func (r *Room) Leave(ch chan string) {

}

var room = Room{
	members: make([]chan string, 0, 0),
}

func hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ch := make(chan string)
	room.Join(ch)

	go func() {
		for msg := range ch {
			// Write
			err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				room.Leave(ch)
				c.Logger().Error(err)
			}
		}
	}()

	for {
		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			room.Leave(ch)
			c.Logger().Error(err)
		}
		room.Send(msg)
		fmt.Printf("%s\n", msg)
	}

}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./public")
	e.GET("/ws", hello)
	e.Logger.Fatal(e.Start(":1323"))
}
