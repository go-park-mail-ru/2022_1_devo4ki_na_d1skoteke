package entity

import "mime/multipart"

type ImageUnit struct {
	Payload     multipart.File
	PayloadSize int64
}
