# ScreenShot

A small HTTP server that takes a screenshot of a web page.

using [Chromium](https://github.com/chromedp/chromedp) to capture the screenshot and [BBolt](https://go.etcd.io/bbolt) to cache images.

# Install & Usage
 - Download
```
git clone https://github.com/issmeftah/screenshot
```
 - Install
```
go build -o bin/screenshot -v main.go
```
- Config
app configuration
```yaml
# port to start HTTP server.
port: "81"
# show chromium debug
debug: false
# http request rate limit: 1 req/sec and permits 2 bursts
rate: 1
bursts: 2
# how long images will be stored
image_cache: 10m
# paths
store: "data/app_store"
log: "data/app_log"
chrome_data: "data/chrome"
```
```
./bin/screenshot run -c "config.yaml"
```

## Capture
 - HTTP `POST /capture`
```json
{
    "url":"https://google.com",
}
```

 - Response example
```json
{
    "Status": 200,
    "Payload": {
        "uuid": "c3b5ohg6n88rm0g59b30"
    }
}
```

## Download
 - HTTP `GET /download/{uuid}`

 if the image was captured successfully this will return the captured image.

 in case there was an error capturing the images or the images is still pending captured, a message will be return instead of the emails

 - Response example 
 ```json
 {
    "Status": 500,
    "Payload": "page load error net::ERR_NAME_NOT_RESOLVED"
}
 ```
 or 
 ```json
 {
    "Status": 500,
    "Payload": "screenshot not yet captured."
}
 ```

## Stats
 - HTTP `POST /stats`

this will return general stats and resource usage.

- Response example
```json
{
    "Status": 200,
    "Payload": {
        "Limiter": {
            "total": 1
        },
        "Resources": {
            "CPUs": 4,
            "CompletedGC": 5,
            "GCSize": "4.47M",
            "Goroutine": 4,
            "StackSystem": "480K",
            "System": "70.83M"
        },
        "Store": {
            "fail": 0,
            "pending": 0,
            "success": 1
        }
    }
}
```