package main

import (
	"log"
	"github.com/sacOO7/gowebsocket"
	"os"
	"os/signal"
	"encoding/json"
	"time"
	"net"
	"github.com/tdesaules/snapwatcher/lib"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ws_url := "ws://" + os.Getenv("SNAPSERVER_HOST") + ":" + os.Getenv("SNAPSERVER_PORT") + "/jsonrpc"
	ws := gowebsocket.New(ws_url)

	ws.OnConnected = func(ws gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	ws.OnConnectError = func(err error, ws gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	ws.OnTextMessage = func(message string, ws gowebsocket.Socket) {
		log.Println("--------------------------------------------------------------")
		log.Println("Recieved message " + message)
		log.Println("--------------------------------------------------------------")
		var data snapcast.Snapcast
		json.Unmarshal([]byte(message), &data)
		if (data.Method == "Stream.OnUpdate") && (data.Params.Stream.Status == "playing") {
			// time.Sleep(0 * time.Millisecond)
			// snapcast.ConfigureStream(data.Params.Stream.Id, data.Params.Stream.Status)
		}
		if data.Method == "Client.OnConnect" {
			snapcast.ConfigureGroupName(data.Params.Client.Config.Name)
		}
		// if data.Id == 69443529 {
		// 	var rooms []snapcast.Room
		// 	var room string
		// 	var group_id string
		// 	json.Unmarshal([]byte(os.Getenv("MULTIROOM")), &rooms)
		// 	for stream_index := range data.Result.Server.Streams {
		// 		group_id = ""
		// 		if data.Result.Server.Streams[stream_index].Id != "idle" {
		// 			room = strings.Split(data.Result.Server.Streams[stream_index].Id, "_")[1]
		// 		}
		// 		if data.Result.Server.Streams[stream_index].Status == "playing" {
		// 			for group_index := range data.Result.Server.Groups {
		// 				if (group_id == "") && (len(data.Result.Server.Groups[group_index].Clients) == 1) && (data.Result.Server.Groups[group_index].Clients[0].Config.Name == room) {
		// 					log.Println("Case 1")
		// 					group_id = data.Result.Server.Groups[group_index].Id
		// 					snapcast.ConfigureGroup(data, group_id, data.Result.Server.Streams[stream_index].Id, room)
		// 				}
		// 				if (group_id == "") && (data.Result.Server.Streams[stream_index].Id == data.Result.Server.Groups[group_index].StreamId) {
		// 					log.Println("Case 2")
		// 					group_id = data.Result.Server.Groups[group_index].Id
		// 					snapcast.ConfigureGroup(data, group_id, data.Result.Server.Streams[stream_index].Id, room)
		// 				}
		// 				if (group_id == "") && (snapcast.GetStreamStatusFromId(data, data.Result.Server.Groups[group_index].StreamId) == "idle") {
		// 					log.Println("Case 3")
		// 					group_id = data.Result.Server.Groups[group_index].Id
		// 					snapcast.ConfigureGroup(data, group_id, data.Result.Server.Streams[stream_index].Id, room)
		// 				}
		// 			}
		// 		}
		// 	}
		// }
	}

	ws.OnDisconnected = func(err error, ws gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		up := false
		timeout := 5 * time.Second
		for up == false {
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(os.Getenv("SNAPSERVER_HOST"), os.Getenv("SNAPSERVER_PORT")), timeout)
			if err != nil {
				log.Println(err)
			}
			if conn != nil {
				defer conn.Close()
				up = true
				log.Println("server is up")
				ws.Connect()
			}
		}
		return
	}

	ws.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			ws.Close()
			return
		}
	}
}