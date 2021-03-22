// Utility to proxy http requests into telegram messages to specific chat via API token
package main

import (
    "fmt"
    "log"
    "flag"
    "net/http"
    "encoding/json"
)

type Response struct {
    Status string `json:"status"`
    Data string `json:"data"`
}

type Payload struct {
    Code string `json:"code"`
    Msg string `json:"msg"`
}

// Generte function to handle HTTP requests and send notifiications
// to telegram
func generate_notification_processor(token, chatid, code string)  func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        payload := Payload{
            Code: query.Get("code"),
            Msg: query.Get("msg"),
        }
        if payload.Code == code {
            send_notification(token, chatid, payload.Msg)
        }
        response := Response{
            Status: "ok",
        }
        json.NewEncoder(w).Encode(response)
    }
}

func send_notification(token, chatid, message string) {
    var url string
    url = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=\"%s\"", token, chatid, message)
    _, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
    }
    // return resp, err
}

func main() {
    var token string
    var chatid string
    var code string
    var bind string

    flag.StringVar(&token, "token", "", "Telegram bot token")
    flag.StringVar(&chatid, "chat", "", "Telegram ChatID")
    flag.StringVar(&code, "code", "42", "Quasy seecret code to match in get request")
    flag.StringVar(&bind, "port", ":8000", "Address/port to bind to")
    flag.Parse()
    if token == "" {
        log.Fatal("token should be defined")
    }
    if chatid == "" {
        log.Fatal("chatid should be defined")
    }
    http.HandleFunc("/", generate_notification_processor(token, chatid, code))
    send_notification(token, chatid, "listening")
    http.ListenAndServe(bind, nil)
}
