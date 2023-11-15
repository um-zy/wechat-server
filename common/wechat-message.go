package common

import (
    "encoding/json"
    "fmt"
    "log"
    "encoding/xml"
    "github.com/go-resty/resty/v2"
)
type WeChatMessageRequest struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
	MsgDataId    int64    `xml:"MsgDataId"`
	Idx          int64    `xml:"Idx"`
}

type WeChatMessageResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}


const (
    apiEndpoint = "https://openai.littlezy.top"
)


func ProcessWeChatMessage(req *WeChatMessageRequest, res *WeChatMessageResponse) {
	switch req.Content {
	case "验证码":
		code := GenerateAllNumberVerificationCode(6)
		RegisterWeChatCodeAndID(code, req.FromUserName)
		res.Content = code
	}
	apiKey := "sk-cypO9JshihUK7NvP36643dD411Ad4fAaB1D37fE121700962"
    client := resty.New()

    response, err := client.R().
        SetAuthToken(apiKey).
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "model":      "ERNIE-Bot",
            "messages":   []interface{}{map[string]interface{}{"role": "system", "content": req.Content}},
            "max_tokens": 50,
        }).
        Post(apiEndpoint)

    if err != nil {
        log.Fatalf("Error while sending send the request: %v", err)
    }

    body := response.Body()

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("Error while decoding JSON response:", err)
        return
    }

    // Extract the content from the JSON response
    content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
    fmt.Println(content)
	res.Content=content

}
