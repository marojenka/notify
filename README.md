# HTTP to TG notification

Simple util to proxy http requests to telegram channel. 
Start server with:
`notify -p :8000 --token $TOKEN --chat $CHATID --code secretcode`

Where:
+ TOKEN telegram api token
+ CHATID id of channel where to send message
+ secretcode arbitrary code to be used as primitive authentication

## Send messages
message can be send via get queryset parameter or post request
```
curl http://127.0.0.1:8000/?code=secretcode&msg=boob
curl http://127.0.0.1:8000/secretcode/ --data-binary "boop"
```

## Notice
proxied post request will be limited by arbitrary length.  Use on your own
risk. IE don't use it. Why would you?
