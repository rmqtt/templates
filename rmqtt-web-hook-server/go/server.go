package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

//client_connect
type ClientConnect struct {
	Action       string `mapstructure:"action"`
	Node         uint64 `mapstructure:"node"`
	Ipaddress    string `mapstructure:"ipaddress"`
	Clientid     string `mapstructure:"clientid"`
	Username     string `mapstructure:"username"`
	Keepalive    uint16 `mapstructure:"keepalive"`
	ProtoVer     int8   `mapstructure:"proto_ver"`
	CleanSession bool   `mapstructure:"clean_session"` //MQTT 3.1, 3.1.1
	CleanStart   bool   `mapstructure:"clean_start"`   //MQTT 5.0
}

//client_connected
type ClientConnected struct {
	Action       string `mapstructure:"action"`
	Node         uint64 `mapstructure:"node"`
	Ipaddress    string `mapstructure:"ipaddress"`
	Clientid     string `mapstructure:"clientid"`
	Username     string `mapstructure:"username"`
	Keepalive    uint16 `mapstructure:"keepalive"`
	ProtoVer     int8   `mapstructure:"proto_ver"`
	CleanSession bool   `mapstructure:"clean_session"` //MQTT 3.1, 3.1.1
	CleanStart   bool   `mapstructure:"clean_start"`   //MQTT 5.0

	ConnectedAt    int64 `json:"connected_at"`
	SessionPresent bool  `json:"session_present"`
}

//client_disconnected
type ClientDisconnected struct {
	Action         string `mapstructure:"action"`
	Node           uint64 `mapstructure:"node"`
	Ipaddress      string `mapstructure:"ipaddress"`
	Clientid       string `mapstructure:"clientid"`
	Username       string `mapstructure:"username"`
	DisconnectedAt int64  `mapstructure:"disconnected_at"`
	Reason         string `mapstructure:"reason"`
}

//message_publish
type MessagePublish struct {
	Action   string `mapstructure:"action"`
	Dup      bool   `mapstructure:"dup"`
	Retain   bool   `mapstructure:"retain"`
	Qos      uint8  `mapstructure:"qos"`
	Topic    string `mapstructure:"topic"`
	Packetid uint16 `mapstructure:"packet_id"`
	Payload  string `mapstructure:"payload"`
	Ts       int64  `mapstructure:"ts"`
}

//message_delivered
type MessageDelivered struct {
	Action   string `mapstructure:"action"`
	Dup      bool   `mapstructure:"dup"`
	Retain   bool   `mapstructure:"retain"`
	Qos      uint8  `mapstructure:"qos"`
	Topic    string `mapstructure:"topic"`
	Packetid uint16 `mapstructure:"packet_id"`
	Payload  string `mapstructure:"payload"`
	Pts      int64  `mapstructure:"pts"`
	Ts       int64  `mapstructure:"ts"`
}

//message_acked
type MessageAcked struct {
	Action   string `mapstructure:"action"`
	Dup      bool   `mapstructure:"dup"`
	Retain   bool   `mapstructure:"retain"`
	Qos      uint8  `mapstructure:"qos"`
	Topic    string `mapstructure:"topic"`
	Packetid uint16 `mapstructure:"packet_id"`
	Payload  string `mapstructure:"payload"`
	Pts      int64  `mapstructure:"pts"`
	Ts       int64  `mapstructure:"ts"`
}

//message_dropped
type MessageDropped struct {
	Action   string `mapstructure:"action"`
	Dup      bool   `mapstructure:"dup"`
	Retain   bool   `mapstructure:"retain"`
	Qos      uint8  `mapstructure:"qos"`
	Topic    string `mapstructure:"topic"`
	Packetid uint16 `mapstructure:"packet_id"`
	Payload  string `mapstructure:"payload"`
	Reason   string `mapstructure:"reason"`
	Pts      int64  `mapstructure:"pts"`
	Ts       int64  `mapstructure:"ts"`
}

func main() {
	router := gin.Default()

	router.POST("/mqtt/webhook", func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			params := make(map[string]interface{})
			if c.Request.Body != nil {
				data, err := ioutil.ReadAll(c.Request.Body)
				if err != nil {
					c.String(http.StatusBadRequest, err.Error())
					return
				}

				err = json.Unmarshal(data, &params)
				if err != nil {
					c.String(http.StatusBadRequest, err.Error())
					return
				}
				fmt.Println(params)
				action := params["action"]
				switch action {
				case "client_connect":
					var client_connect ClientConnect
					err = mapstructure.Decode(params, &client_connect)
					if err != nil {
						fmt.Println(err.Error())
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("*** client_connect: ", client_connect)
				case "client_connected":
					var client_connected ClientConnected
					err = mapstructure.Decode(params, &client_connected)
					if err != nil {
						fmt.Println(err.Error())
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("*** client_connected: ", client_connected)
				case "client_disconnected":
					var client_disconnected ClientDisconnected
					err = mapstructure.Decode(params, &client_disconnected)
					if err != nil {
						fmt.Println(err.Error())
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("*** client_disconnected: ", client_disconnected)

				case "message_publish":
					var message_publish MessagePublish
					err = mapstructure.Decode(params, &message_publish)
					if err != nil {
						fmt.Println(err.Error())
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("*** message_publish: ", message_publish)
					payload_bytes, err1 := base64.StdEncoding.DecodeString(message_publish.Payload)
					if err1 != nil {
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("payload_bytes: ", string(payload_bytes))

				case "message_delivered":
					var message_delivered MessageDelivered
					err = mapstructure.Decode(params, &message_delivered)
					if err != nil {
						fmt.Println(err.Error())
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("*** message_delivered: ", message_delivered)

				case "message_acked":
					var message_acked MessageAcked
					err = mapstructure.Decode(params, &message_acked)
					if err != nil {
						fmt.Println(err.Error())
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("*** message_acked: ", message_acked)

				case "message_dropped":
					var message_dropped MessageDropped
					err = mapstructure.Decode(params, &message_dropped)
					if err != nil {
						fmt.Println(err.Error())
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					fmt.Println("*** message_dropped: ", message_dropped)
				}
			}
		} else {
			c.String(http.StatusBadRequest, fmt.Sprintf("content type(%s) not supported", c.ContentType()))
			return
		}
	})

	router.Run(":5656")
}
