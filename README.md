# elastic-image
 * this is an implementation of the test task
 * starts http server
 * implements JSON API
 * downloads image by URL
   * request with duplicated URL will wait for download ordered by the first request and will use same fetched image
   * when image is broadcasted to all subscribers, it will be wiped
 * resize it with width and heigth passed into the request
 * returns resized image to the client

## Launch linter with command
```bash
golangci-lint run
```

## Tests launch command
```bash
 go test -count=1 -cover -gcflags "all=-l" ./...
```
## Usage example
```bash
curl -s -X POST -H 'Content-Type: application/json' -d \
  '{
     "url":"https://i.pinimg.com/originals/e8/cc/62/e8cc621cc2dda8b2aae42e59140c12ad.jpg",
     "width":1024,
     "heigth":800
   }' \
  http://localhost:8080/api/v1/images/resize \
  -o out.jpg
```
