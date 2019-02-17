# Generating mock
```
cd <PROJECT_ROOT>
mockery -dir PATH_TO_image.go  -output ./internal/tests/mocks -name=Image
```

# After generation
import should be replaced manually:
```
import image "path_to_package_image"
```
to
```
import image "image"
```