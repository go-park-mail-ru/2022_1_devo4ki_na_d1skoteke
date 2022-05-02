package xss

import (
	"cotion/internal/domain/entity"
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

func SanitizeNotes(data *entity.ShortNotes) {
	if sanitizer == nil {
		return
	}
	for i := 0; i < len((*data).ShortNote); i++ {
		(*data).ShortNote[i].Name = sanitizer.Sanitize((*data).ShortNote[i].Name)
		(*data).ShortNote[i].Body = sanitizer.Sanitize((*data).ShortNote[i].Body)
		(*data).ShortNote[i].Token = sanitizer.Sanitize((*data).ShortNote[i].Token)
	}
}

func SanitizeNote(data *entity.Note) {
	if sanitizer == nil {
		return
	}
	(*data).Name = sanitizer.Sanitize((*data).Name)
	(*data).Body = sanitizer.Sanitize((*data).Body)
}
