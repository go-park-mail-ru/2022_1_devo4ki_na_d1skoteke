package xss

import (
	"cotion/internal/api/domain/entity"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

var sanitizer *bluemonday.Policy = nil

func NewXssSanitizer() {
	sanitizer = bluemonday.UGCPolicy()
	log.Info("Xss sanitizer is on")
}

func Sanitize(rawText string) string {
	if sanitizer == nil {
		return rawText
	}
	return sanitizer.Sanitize(rawText)
}

func SanitizeNotes(data *[]entity.Note) {
	if sanitizer == nil {
		return
	}
	for i := 0; i < len(*data); i++ {
		(*data)[i].Name = sanitizer.Sanitize((*data)[i].Name)
		(*data)[i].Body = sanitizer.Sanitize((*data)[i].Body)
	}
}

func SanitizeNote(data *entity.Note) {
	if sanitizer == nil {
		return
	}
	(*data).Name = sanitizer.Sanitize((*data).Name)
	(*data).Body = sanitizer.Sanitize((*data).Body)
}
