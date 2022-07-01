package main

import (
  "fmt"
  "os"
  "sync"

  MQTT "github.com/eclipse/paho.mqtt.golang"
  "github.com/gorilla/websocket"
)

type Message struct {
  Jsonrpc string                 `json:"jsonrpc"`
  Method  string                 `json:"method"`
  Params  map[string]interface{} `json:"params"`
}

var ws *websocket.Conn
var mqttClient MQTT.Client

func onWebSocketMessage(message []byte) {
  // fmt.Println(string(message))
  mqttClient.Publish("snapcast/rx", 0, false, string(message))
}

func onMqttMessage(client MQTT.Client, msg MQTT.Message) {
  // fmt.Printf("%s\n", msg.Payload())
  err := ws.WriteMessage(websocket.TextMessage, msg.Payload())
  if err != nil {
    fmt.Printf("[ERROR] Sending WebSocket message:", err)
  }
}

func initWebSocketClient() {
  snapcast_ip := os.Getenv("SNAPCAST_IP")

  var err error

  ws, _, err = websocket.DefaultDialer.Dial("ws://"+snapcast_ip+":1780/jsonrpc", nil)
  if err != nil {
    fmt.Printf("[ERROR] Connecting to WebSocket:", err)
    os.Exit(1)
  }
  fmt.Println("Connected to WebSocket")

  go func() {
    for {
      _, message, err := ws.ReadMessage()
      if err != nil {
        fmt.Println("[ERROR] Reading WebSocket message:", err)
        os.Exit(1)
      }
      onWebSocketMessage(message)
    }
  }()
  fmt.Println("Listening on WebSocket")
}

func initMqttClient() {
  mqtt_ip := os.Getenv("MQTT_IP")

  opts := MQTT.NewClientOptions().AddBroker("tcp://" + mqtt_ip + ":1883").SetClientID("snapcast")
  opts.SetCleanSession(true)

  // Connect to MQTT
  mqttClient = MQTT.NewClient(opts)
  if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
    fmt.Println("[ERROR] Connecting to MQTT:", token.Error())
    os.Exit(1)
  }
  fmt.Println("Connected to MQTT")

  // MQTT Listener
  if token := mqttClient.Subscribe("snapcast/tx", 0, onMqttMessage); token.Wait() && token.Error() != nil {
    fmt.Println("[ERROR] Reading MQTT message:", token.Error())
    os.Exit(1)
  }
  fmt.Println("Listening on MQTT")
}

func waitForever() {
  var wg sync.WaitGroup
  wg.Add(1)
  wg.Wait()
}

func main() {
  initWebSocketClient()
  initMqttClient()
  waitForever()
}
