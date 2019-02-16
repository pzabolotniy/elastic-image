package entity

type ImageInfo struct {
	URL string  `json:"url" validate:"required"`
	Width int32 `json:"width" validate:"required"`
	Height int32 `json:"heigth" validate:"required"`
}
