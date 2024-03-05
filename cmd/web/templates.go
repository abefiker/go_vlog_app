package main

import (
	"github.com/abefiker/go_vlog_app/internal/models"
)

type templateData struct {
	Vlog *models.Vlog
	Vlogs []*models.Vlog
}
