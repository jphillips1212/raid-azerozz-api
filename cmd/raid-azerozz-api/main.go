package main

import (
	api "github.com/jphillips1212/roztools-api/api"
	handlers "github.com/jphillips1212/roztools-api/api/handlers"
	fs "github.com/jphillips1212/roztools-api/internal/db/firestore"
	wl "github.com/jphillips1212/roztools-api/internal/logs/warcraftlogs"
)

func main() {
	db := fs.New()
	wl := wl.New()

	warcraft := handlers.Handler{
		Database: db,
		Logs:     wl,
	}

	api.Register(warcraft)
}
