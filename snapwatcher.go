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
	"strings"
)

type Room struct {
	Name string `json:"name"`
	Members []string `json:"members"`
}

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
			time.Sleep(1500 * time.Millisecond)
			ws.SendText(`{"id":69443529,"jsonrpc":"2.0","method":"Server.GetStatus"}`)
		}
		if data.Id == 69443529 {
			var rooms []Room
			var request string
			var group_id string
			var clients_id []string
			json.Unmarshal([]byte(os.Getenv("MULTIROOM")), &rooms)
			for stream_index := range data.Result.Server.Streams {
				if data.Result.Server.Streams[stream_index].Status == "playing" {
					for group_index := range data.Result.Server.Groups {
						if snapcast.GetStreamStatusFromId(data, data.Result.Server.Groups[group_index].StreamId) == "idle" {
							group_id = data.Result.Server.Groups[group_index].Id
						} else {
							log.Println("stream are already playing sound")
						}
					}
					if group_id != "" {
						request = `{"id":12931886,"jsonrpc":"2.0","method":"Group.SetStream","params":{"id":"` + group_id + `","stream_id":"` + data.Result.Server.Streams[stream_index].Id + `"}}`
						ws.SendText(request)
						log.Println("set stream " + data.Result.Server.Streams[stream_index].Id + " on group " + group_id )
						request = `{"id":18029639,"jsonrpc":"2.0","method":"Group.SetName","params":{"id":"` + group_id + `","name":"` + data.Result.Server.Streams[stream_index].Id + `"}}`
						ws.SendText(request)
						log.Println("set groupe name " + data.Result.Server.Streams[stream_index].Id + " on group " + group_id )
					}
					for room_index := range rooms {
						room := strings.Split(data.Result.Server.Streams[stream_index].Id, "_")[1]
						clients_id = nil
						if room == rooms[room_index].Name {
							for member_index := range rooms[room_index].Members {
								client_id := snapcast.GetClientIdFromName(data, rooms[room_index].Members[member_index])
								clients_id = append(clients_id, client_id)
							}
							if clients_id != nil {
								json_clients_id, _ := json.Marshal(clients_id)
								request = `{"id":85139337,"jsonrpc":"2.0","method":"Group.SetClients","params":{"clients":` + string(json_clients_id) + `,"id":"` + group_id +`"}}`
								if len(clients_id) == 1 {
									ws.SendText(request)
								} else {
									ws.SendText(request)
								}
							}
						}

					}
				}
			}
		}
		if data.Id == 85139337 {
			for group_index := range data.Result.Server.Groups {
				if data.Result.Server.Groups[group_index].Name != data.Result.Server.Groups[group_index].StreamId {
					log.Println("orphan group")
				}
			}
		}
		if data.Id == 12931886 {
			log.Println("Stream id " + data.Result.StreamId + " has been apply")
		}
		if data.Id == 18029639 {
			log.Println("Group name " + data.Result.Name + " has been apply")
		}
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