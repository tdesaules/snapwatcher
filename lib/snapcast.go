package snapcast

import (
	"github.com/ybbus/jsonrpc/v2"
	"os"
	// "encoding/json"
	// "strings"
)

type Config struct {
	Name string `json:"name"`
}
type Client struct {
	Id string `json:"id"`
	Config Config `json:"config"`
}
type Group struct {
	Clients []Client `json:"clients"`
	Id string `json:"id"`
	Name string `json:"name"`
	StreamId string `json:"stream_id"`
}
type Stream struct {
	Id string `json:"id"`
	Status string `json:"status"`
}
type Server struct {
	Groups []Group `json:"groups"`
	Streams []Stream `json:"streams"`
}
type Result struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
	StreamId string `json:"stream_id"`
	Name string `json:"name"`
	Server Server `json:"server"`
}
type Params struct {
	Id string `json:"id"`
	Stream Stream `json:"stream"`
	Client Client `json:"client"`
}
type Error struct {
	Code int `json:"code"`
	Data string `json:"data"`
	Message string `json:"message"`
}
type Snapcast struct {
	Id int `json:"id"`
	Result Result `json:"result"`
	Method string `json:"method"`
	Params Params `json:"params"`
	Error Error `json:"error"` 
}
type GroupSetStream struct {
	Id string `json:"id"`
	StreamId string `json:"stream_id"`
}
type GroupSetName struct {
	Id string `json:"id"`
	Name string `json:"name"`
}
type GroupSetClients struct {
	Id string `json:"id"`
	Clients []string `json:"clients"`
}
type Room struct {
	Name string `json:"name"`
	Members []string `json:"members"`
}

// func GetStreamStatusFromId(data Result, id string) (status string) {
// 	for stream_index := range data.Server.Streams {
// 		if data.Server.Streams[stream_index].Id == id {
// 			status = data.Server.Streams[stream_index].Status
// 		}
// 	}
// 	return status
// }

func GetClientIdFromName(name string) (client_id string) {
	var data Result
	rpc_client := jsonrpc.NewClient("http://" + os.Getenv("SNAPSERVER_HOST") + ":" + os.Getenv("SNAPSERVER_PORT") + "/jsonrpc")
	response, _ := rpc_client.Call("Server.GetStatus")
	response.GetObject(&data)

	for group_index := range data.Server.Groups {
		for client_index := range data.Server.Groups[group_index].Clients {
			if data.Server.Groups[group_index].Clients[client_index].Config.Name == name {
				client_id = data.Server.Groups[group_index].Clients[client_index].Id
			}
		}
	}

	return client_id
}

// func GetClientsNameForRoom(room string) (clients_name []string) {
// 	var rooms []Room
// 	json.Unmarshal([]byte(os.Getenv("MULTIROOM")), &rooms)
// 	for room_index := range rooms {
// 		if room == rooms[room_index].Name {
// 			for member_index := range rooms[room_index].Members {
// 				client_name := rooms[room_index].Members[member_index]
// 				clients_name = append(clients_name, client_name)
// 			}
// 		}
// 	}
// 	return clients_name
// }

// func ConfigureGroup(data Result, group_id string, stream_id string, room string) {
// 	rpc_client := jsonrpc.NewClient("http://" + os.Getenv("SNAPSERVER_HOST") + ":" + os.Getenv("SNAPSERVER_PORT") + "/jsonrpc")
// 	rpc_client.Call("Group.SetStream", &GroupSetStream{group_id, stream_id})
// 	rpc_client.Call("Group.SetName", &GroupSetName{group_id, stream_id})
// 	if stream_id != "idle" {
// 		rpc_client.Call("Group.SetClients", &GroupSetClients{group_id, GetClientsNameForRoom(data, room)})
// 	}
// }

// func ConfigureStream(stream_id string, stream_status string) {
// 	if stream_status == "playing" {
// 		var data Result
// 		var group_id string
// 		var room string

// 		rpc_client := jsonrpc.NewClient("http://" + os.Getenv("SNAPSERVER_HOST") + ":" + os.Getenv("SNAPSERVER_PORT") + "/jsonrpc")
// 		response, _ := rpc_client.Call("Server.GetStatus")
// 		response.GetObject(&data)

// 		if stream_id != "idle" {
// 			room = strings.Split(stream_id, "_")[1]
// 		}

// 		for group_index := range data.Server.Groups {

// 		}
// 	}

// }

func ConfigureGroupName(client_name string) {
	var data Result
	rpc_client := jsonrpc.NewClient("http://" + os.Getenv("SNAPSERVER_HOST") + ":" + os.Getenv("SNAPSERVER_PORT") + "/jsonrpc")
	response, _ := rpc_client.Call("Server.GetStatus")
	response.GetObject(&data)

	client_id := GetClientIdFromName(client_name)

	for group_index := range data.Server.Groups {
		if len(data.Server.Groups[group_index].Clients) == 1 {
			if data.Server.Groups[group_index].Clients[0].Id == client_id {
				rpc_client.Call("Group.SetName", &GroupSetName{data.Server.Groups[group_index].Id, client_name})
			}
		}
	}
}