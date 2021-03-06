// Utility to proxy http requests into telegram messages to specific chat via API token
package main

import (
    "io"
    "fmt"
    "log"
    "flag"
    "net/url"
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

// Generate function to extract message from get parameters and send
// notification
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

// Generate function to extract message from request body and send it as notification
func generate_body_processor(token, chatid string)  func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
	    limited_body := io.LimitReader(r.Body, 1000)
        msg, err := io.ReadAll(limited_body)
        if err != nil {
            log.Printf("Error reading body: %v", err)
            http.Error(w, "can't read body", http.StatusBadRequest)
            return
        }
	    fmt.Printf("%s", msg)
        send_notification(token, chatid, fmt.Sprintf("%s", msg))
        response := Response{
            Status: "ok",
        }
        json.NewEncoder(w).Encode(response)
    }
}

func send_notification(token, chatid, message string) {
    req, err := url.Parse(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token))
    if err != nil {
        log.Fatal(err)
    }

    q := url.Values{}
    q.Add("chat_id", chatid)
    q.Add("text", message)
    req.RawQuery = q.Encode()
    _, errr := http.Get(req.String())
    if errr != nil {
        log.Printf("Error sending notification: %v", errr)
    }
    fmt.Printf("%s", message)
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
    http.HandleFunc(fmt.Sprintf("/%s/", code), generate_body_processor(token, chatid))
    send_notification(token, chatid, "listening")
    http.ListenAndServe(bind, nil)
}
