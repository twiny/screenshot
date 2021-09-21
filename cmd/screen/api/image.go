package api

import (
	"time"
)

// ImageStatus
type ImageStatus string

const (
	ImageStatusSuccess ImageStatus = "success"
	ImageStatusPending ImageStatus = "pending"
	ImageStatusFail    ImageStatus = "fail"
)

// Image
type Image struct {
	UUID      string      `json:"uuid"`
	Status    ImageStatus `json:"status"`
	Message   string      `json:"message"`
	Binary    []byte      `json:"binary"`
	CreatedAt time.Time   `json:"created_at"`
}
