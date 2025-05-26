// @title           Quiz API
// @version         1.0
// @description     API Server for Quiz Application
// @host            localhost:8080
// @BasePath        /
package main

import (
	"learn_phase_2_local_server/db"
	"learn_phase_2_local_server/router"
)

func main() {
	db.Init()
	r := router.SetupRouter()
	r.Run()
}
