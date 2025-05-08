package canvas

import (
	"time"
)

type CanvasBase struct {
	CanvasID  *int      `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      any       `json:"data"`
	Text      string    `json:"text"`
	FolderId  int       `json:"folder_id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CanvasList struct {
	Canvases []*CanvasBase `json:"canvases"`
}

type FolderBase struct {
	ID        *int      `json:"id"`
	Name      string    `json:"name"`
	Parent    int       `json:"parent_folder_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type FolderList struct {
	Folders []*FolderBase `json:"folders"`
}

type FolderContent struct {
	Folder      FolderBase   `json:"folder"`
	BreadCrumbs []BreadCrumb `json:"breadcrumbs"`
	Folders     []FolderBase `json:"folders"`
	Canvases    []CanvasBase `json:"canvases"`
}

type BreadCrumb struct {
	ID   *int   `json:"id"`
	Name string `json:"name"`
}

// type BreadCrumbs struct {
// 	Crumbs []BreadCrumbs `json:"breadcrumbs"`
// }
