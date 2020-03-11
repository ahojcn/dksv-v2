package models

import "time"

type ImageInfo struct {
	Name    string      `json:"name"`
	Sys     interface{} `json:"sys"`
	ModTime time.Time   `json:"mod_time"`
	Size    int64       `json:"size"`
}
