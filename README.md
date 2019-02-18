# elastic-image
image resize service

## Tests launch command
```
go test -gcflags "all=-l" ./...
```
## Usage example
```bash
curl -X POST -H 'Content-Type: application/json' -d '{"url":"https://cdn.wallpapersbuzz.com/image/1948/b_mountains-view.jpg","width":1024,"heigth":800}' http://localhost:8080/api/v1/images/resize
```
