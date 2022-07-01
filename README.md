# snapcast2mqtt

A bridge between snapcast's websocket and an MQTT instance

It's basically forwarding between WebSocket and MQTT all the messages without any processing. This allows easy Snapcast control over MQTT.

### Build

```sh
go build
```

### Requirements

You'll need to set the following env variables which point to the Snapcast & MQTT servers

```sh
SNAPCAST_IP=192.168.0.xxx
MQTT_IP=192.168.0.xxx
```

### Run

```sh
./snapcast2mqtt
```
