package snapcast

import (

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
}

type Error struct {
	Code int `json:"code"`
	Data string `json:"data"`
	Message string `json:"message"`
}

type Snapcast struct {
	Id int `json:"id"`
	Jsonrpc string `json:"jsonrpc`
	Result Result `json:"result"`
	Method string `json:"method"`
	Params Params `json:"params"`
	Error Error `json:"error"` 
}

func GetStreamStatusFromId(data Snapcast, id string) (status string) {
	for stream_index := range data.Result.Server.Streams {
		if data.Result.Server.Streams[stream_index].Id == id {
			status = data.Result.Server.Streams[stream_index].Status
		}
	}
	return status
}

func GetClientIdFromName(data Snapcast, name string) (client_id string) {
	for group_index := range data.Result.Server.Groups {
		for client_index := range data.Result.Server.Groups[group_index].Clients {
			if data.Result.Server.Groups[group_index].Clients[client_index].Config.Name == name {
				client_id = data.Result.Server.Groups[group_index].Clients[client_index].Id
			}
		}
	}
	return client_id
}