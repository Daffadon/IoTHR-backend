package main

import (
	"IoTHR-backend/db"
	"IoTHR-backend/migrations"
	"IoTHR-backend/server"
)

func main() {
	db.Init()
	migrations.Init()
	server.Init()
}
