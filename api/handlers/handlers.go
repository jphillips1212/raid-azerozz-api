package handlers

import (
	db "github.com/jphillips1212/roztools-api/internal/db"
	logs "github.com/jphillips1212/roztools-api/internal/logs"
	svc "github.com/jphillips1212/roztools-api/internal/service"
)

type Handler struct {
	Database db.DB
	Logs     logs.WarcraftLogs
}

type Service struct {
	Service svc.Service
}
