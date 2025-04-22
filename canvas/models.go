package canvas

import (
	"time"
)

type CanvasBase struct {
	CanvasID  *int      `json:"id"`
	Name      string    `json:"canvas_name"`
	Url       string    `json:"canvas_url"`
	UpdatedAt time.Time `json:"update_time"`
	Data      any       `json:"data"`
	Text      string    `json:"text"`
}

type CanvasList struct {
	Canvases []*CanvasBase `json:"canvases"`
}
