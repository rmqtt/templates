package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
const CONTENT_TYPE_JSON = "application/json"

func auth(c *gin.Context) {
	params, err := parseParams(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		clientid := params.Get("clientid")
		username := params.Get("username")
		password := params.Get("password")
		protocol := params.Get("protocol")
		fmt.Println("auth clientid:", clientid)
		fmt.Println("auth username:", username)
		fmt.Println("auth password:", password)
		fmt.Println("auth protocol:", protocol)

		// @TODO Verify user validity,

		// is user-ignore
		if username == "user-ignore" {
			c.String(http.StatusOK, "ignore")
			return
		}

		// is user-deny
		if username == "user-deny" {
			c.String(http.StatusOK, "deny")
			return
		}

		// is admin
		if username == "user-admin" {
			c.Header("X-Superuser", "true")
		}

		//acl
		if username == "user-acl" {
			json_acl := `
			{
				"result":"allow",
				"superuser": false,
				"expire_at": 1827143027,
				"acl": [
					{
					"permission": "allow",
					"action": "all",
					"topic": "foo/${clientid}"
					},
					{
					"permission": "allow",
					"action": "subscribe",
					"topic": "eq foo/1/#",
					"qos": [1,2]
					},
					{
					"permission": "allow",
					"action": "subscribe",
					"topic": "foo/2/#",
					"qos": 1
					},
					{
					"permission": "allow",
					"action": "publish",
					"topic": "foo/2/1",
					"qos": 1
					},
					{
					"permission": "allow",
					"action": "publish",
					"topic": "foo/${username}",
					"retain": false,
					"qos": [0,1]
					},
					{
					"permission": "deny",
					"action": "all",
					"topic": "foo/3"
					},
					{
					"permission": "deny",
					"action": "publish",
					"topic": "foo/4",
					"retain": true
					}
				]
			}`
			var json_data map[string]interface{}
			json.Unmarshal([]byte(json_acl), &json_data)
			c.JSON(http.StatusOK, json_data)
			return
		}

		// allow
		c.String(http.StatusOK, "allow")
	}
}

func acl(c *gin.Context) {
	params, err := parseParams(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		//access = "%A", username = "%u", protocol = "%r", clientid = "%c", ipaddr = "%a", topic = "%t"
		access := params.Get("access")
		clientid := params.Get("clientid")
		username := params.Get("username")
		protocol := params.Get("protocol")
		ipaddr := params.Get("ipaddr")
		topic := params.Get("topic")

		fmt.Println("acl clientid:", clientid)
		fmt.Println("acl username:", username)
		fmt.Println("acl protocol:", protocol)
		fmt.Println("acl access:", access)
		fmt.Println("acl ipaddr:", ipaddr)
		fmt.Println("acl topic:", topic)

		// @TODO Verify topic validity,

		// is Subscribe
		if access == "1" {
			fmt.Println("is Subscribe, topic is ", topic)
		}

		if access == "2" {
			fmt.Println("is Publish, topic is ", topic)
		}

		if strings.HasSuffix(topic, "/cache") {
			c.Header("X-Cache", "-1") //Unit millisecond, if the value is -1, it will not expire before disconnecting
		}

		// test ignore
		if strings.HasPrefix(topic, "/test/ignore") {
			c.String(http.StatusOK, "ignore")
			return
		}

		// test deny
		if strings.HasPrefix(topic, "/test/deny") {
			c.String(http.StatusOK, "deny")
			return
		}

		// allow
		c.String(http.StatusOK, "allow")
	}
}

func main() {
	router := gin.Default()

	router.GET("/mqtt/auth", auth)
	router.POST("/mqtt/auth", auth)
	router.PUT("/mqtt/auth", auth)

	router.GET("/mqtt/acl", acl)
	router.POST("/mqtt/acl", acl)
	router.PUT("/mqtt/acl", acl)

	router.Run(":9090")
}

type Params interface {
	Get(key string) string
}

type QueryParams struct {
	c *gin.Context
}

func (p *QueryParams) Get(key string) string {
	v, _ := p.c.GetQuery(key)
	return v
}

type FormParams struct {
	c *gin.Context
}

func (p *FormParams) Get(key string) string {
	v, _ := p.c.GetPostForm(key)
	return v
}

type JsonParams struct {
	c      *gin.Context
	params map[string]interface{}
}

func (p *JsonParams) parseJson() error {
	if p.params == nil {
		p.params = make(map[string]interface{})
		if p.c.Request.Body != nil {
			data, err := ioutil.ReadAll(p.c.Request.Body)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, &p.params)
		}
	}
	return nil
}

func (p *JsonParams) Get(key string) string {
	if err := p.parseJson(); err != nil {
		log.Println(err)
		return ""
	}
	if v, ok := p.params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func parseParams(c *gin.Context) (Params, error) {
	var params Params
	switch c.Request.Method {
	case "GET":
		params = &QueryParams{c}
	case "POST", "PUT":
		ctype := c.ContentType()
		if ctype == CONTENT_TYPE_FORM {
			params = &FormParams{c}
		} else if ctype == CONTENT_TYPE_JSON {
			params = &JsonParams{c, nil}
		} else {
			return nil, fmt.Errorf("content type(%s) not supported", ctype)
		}
	default:
		return nil, fmt.Errorf("method(%s) not supported", c.Request.Method)
	}
	return params, nil
}
