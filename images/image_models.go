package images

import (
	"time"
)

type (
	Image struct {
		UserId int64  `json:"person_id"`
		Url    string `json:"image_url"`
	}

	Canvas struct {
		Id     int64     `json:"id"`
		Name   string    `json:"canvas_name"`
		Url    string    `json:"canvas_url"`
		Update time.Time `json:"update_time"`
	}

	Canvases struct {
		Canvases []Canvas `json:"canvases"`
	}

	// CanvasFile struct {
	// 	File os.File `json:"file"`
	// }

	CanvasRequest struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
)
