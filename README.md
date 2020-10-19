# elastic-image
 * starts http server
 * implements JSON API
 * downloads image by URL
 * resize it with width and heigth passed into the request
 * returns to the client

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
