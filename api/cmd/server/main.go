package main

import "github.com/RintaroNasu/drive-route-planner/api/routes"

func main() {
	e := routes.NewRouter()

	e.Logger.Fatal(e.Start(":8080"))
}
