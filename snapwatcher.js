#!/usr/bin/env node

const WebSocket = require('ws');
const ws = new WebSocket('ws://192.168.100.173:1780/jsonrpc');

ws.on('message', function incoming(data) {
    data = JSON.parse(data);
    if (data.method == 'Stream.OnUpdate' && data.params.stream.status == 'playing') {
        console.log('stream is playing on: ', data.params.stream.id)
    }
    if (data.method == 'Stream.OnUpdate' && data.params.stream.status == 'idle') {
        console.log('streaming is over on: ', data.params.stream.id)
    }
});