#!/usr/bin/env node

var fetch = require("node-fetch");
var jsonrpc = require('jsonrpc-lite');
var websocket = require('ws');
var opened_socket = false;
var time_interval = 10000;
var multiroom = JSON.parse(process.env.MULTIROOM);
var sources = JSON.parse(process.env.SOURCES);

async function serverGetStatus() {
    var response = await fetch(`http://${process.env.SNAPSERVER_HOST}:${process.env.SNAPSERVER_PORT}/jsonrpc`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
        },
        body: JSON.stringify(jsonrpc.request(32064996, 'Server.GetStatus')),
    });
    var content = await response.json();
    return content;
}
function groupSetStream(ws, group_id, stream_id) {
    ws.send(JSON.stringify(jsonrpc.request(95106597, 'Group.SetStream', {"id":group_id,"stream_id":stream_id})));
}
function groupSetClient(ws, group_id, clients_id) {
    ws.send(JSON.stringify(jsonrpc.request(15548462, 'Group.SetClients', {"clients":clients_id,"id":group_id})));
}
function groupsGetId(data) {
    var groups_id = new Array();
    data.result.server.groups.forEach(function (group) {
        groups_id.push(group.id);
    });
    return groups_id;
}
function groupGetIdForClient(name, data) {
    var group_id;
    var client_id = clientGetId(name, data);
    data.result.server.groups.forEach(function (group) {
        group.clients.filter(function (entry) {
            if (entry.id == client_id) {
                group_id = group.id;
            }
        });
    });
    return group_id;
}
function groupGetIDForStream(name, data) {
    var group_id;
    data.result.server.groups.forEach(function (group) {
        if (group.stream_id == name) {
            group_id = group.id;
        }
    });
    return group_id;
}
function streamGetStatusById(name, data) {
    var stream_status;
    data.result.server.streams.forEach(function (entry) {
        if (entry.id == name) {
            stream_status = entry.status;
        }
    });
    return stream_status;
}
function clientGetId(name, data) {
    var client_id;
    data.result.server.groups.forEach(function (group) {
        var clients = group.clients.filter(function (entry) {
            return entry.config.name == name;
        });
        clients.forEach(function (client) {
            if (client) {
                client_id = client.id;
            }
        })
    });
    return client_id;
}
function streamOnUpdate(ws, data) {
    var source = data.params.id.split("_")[0];
    var room = data.params.id.split("_")[1];
    var stream_id = data.params.id;
    var stream_status = data.params.stream.status;
    var members = multiroom.filter(function (entry) {
        return entry.name == room;
    })[0].members;
    if (stream_status == 'playing') {
        (serverGetStatus()).then(function (result) {
            var clients_id = new Array();
            var group_id;
            members.forEach(function (member) {
                clients_id.push(clientGetId(member, result));
                group_id = groupGetIdForClient(member, result);
            });
            console.log(`stream id : ${stream_id}`);
            console.log(`group id : ${group_id}`);
            console.log(`clients id list : ${clients_id}`);
            console.log('------------------------------------------------------------------');
            groupSetStream(ws, group_id, stream_id);
            groupSetClient(ws, group_id, clients_id);
            // need to add a check if in local stream, check that no other client connected to the group where the stream is connected
        });
    }
    if (stream_status == 'idle') {
        (serverGetStatus()).then(function (result) {
            var group_id = groupGetIDForStream(stream_id, result);
            if (group_id) {
                // get groupid so for that group id get client id to check if they are on the stream and release them
                // placing it in the stream they have to be
            }
        });
    }
}


function connect() {
    var ws = new websocket(`ws://${process.env.SNAPSERVER_HOST}:${process.env.SNAPSERVER_PORT}/jsonrpc`);
    return new Promise((resolve, reject) => {
        ws.on('open', function open() {
            console.log("socket connected");
            opened_socket = true;
            resolve(opened_socket);
        });
        ws.on('message', function incoming(data) {
            //console.log(data);
            data = JSON.parse(data);
            if (data.method == 'Client.OnConnect') {
                console.log(`client ${data.params.client.config.name} connected`);
            }
            if (data.method == 'Client.OnDisconnect') {
                console.log(`client ${data.params.client.config.name} disconnected`);
            }
            if (data.method == 'Stream.OnUpdate') {
                streamOnUpdate(ws, data);
            }
            if (data.id == 11111111) {
                
            }
        });
        ws.on('close', function close(error) {
            console.log('socket disconnected');
            opened_socket = false;
            reject(error);
        });
        ws.on('error', function close(error) {
            console.log('socket error');
            opened_socket = false;
            reject(error);
        });
    });
}
async function reconnect() {
    try {
        await connect();
    } catch {
        console.log("socket disconnected");
    }
}

reconnect();

setInterval(() => {
    if (opened_socket == false) {
        reconnect();
    }
}, time_interval);