package main

import (
	"github.com/juanignaciorc/microbloggin-pltf/cmd/api"
)

func main() {
	router := api.SetupEngine()
	api.StartServer(router)
}
