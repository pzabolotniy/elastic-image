# Generating mock
```
cd <PROJECT_ROOT>
mockery -dir ./internal/httpclient  -output ./internal/httpclient/mocks -name=Browser
mockery -dir ./internal/httpclient  -output ./internal/httpclient/mocks -name=Responser
```