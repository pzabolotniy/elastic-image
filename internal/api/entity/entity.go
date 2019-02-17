package entity

// ImageInfo is a DTO for 'POST /api/v1/images/resize' request
type ImageInfo struct {
	URL    string `json:"url" validate:"required"`
	Width  int32  `json:"width" validate:"required"`
	Height int32  `json:"heigth" validate:"required"`
}
