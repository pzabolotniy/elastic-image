# elastic-image
image resize service

## Tests launch command
```
go test -gcflags "all=-l" ./...
```
## Usage example
```bash
curl -X POST -H 'Content-Type: application/json' -d '{"url":"https://i.pinimg.com/originals/e8/cc/62/e8cc621cc2dda8b2aae42e59140c12ad.jpg","width":1024,"heigth":800}' http://localhost:8080/api/v1/images/resize
```
