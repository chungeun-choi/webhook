package main

import (
	"github.com/chungeun-choi/webhook/service/patch"
	"log"
	"net/http"
)

func main() {
	//if _, err := config.LoadConfig("app.yml"); err != nil {
	//	panic(err)
	//}
	//// Create a client interface implementation here
	//if err := mutating.CreateClientSet(); err != nil {
	//	panic(err)
	//}
	//
	//mutateManager := mutating.NewMutateManager(mutating.ClientCache)
	//
	//// Create the HTTP server and start listening on a port
	//http.ListenAndServe(":8080", mutating.NewMux(mutateManager))

	http.HandleFunc("/addPatch", patch.AddPatchHandler)
	http.HandleFunc("/getPatch", patch.GetPatchHandler)
	http.HandleFunc("/updatePatch", patch.UpdatePatchHandler)
	http.HandleFunc("/clearPatch", patch.ClearPatchHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
