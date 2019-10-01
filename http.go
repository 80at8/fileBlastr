// curl -X POST -H "Content-Type: application/octet-stream" --data-binary '@filename' http://127.0.0.1:5050/upload

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"time"

	"github.com/gorilla/mux"

	"net/http"
)

func HashHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	FileName := vars["filename"]
	ClientHash := vars["shahash"]

	// UUID, ExistingUUID := CheckForExistingUUIDAndHash(ClientHash, FileName)

	// if ExistingUUID == true {
	// 	w.Write([]byte(fmt.Sprintf("Existing UUID %v", UUID)))

	// } else {

	// 	UUID := CreateNewUUIDAndHash(ClientHash)

	// 	w.Write([]byte(fmt.Sprintf("Created UUID %v", UUID)))
	// }
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	FileName := vars["filename"]

	LogRequest(FileName)

	file, err := os.Create("./uploads/" + FileName)
	defer file.Close()

	if err != nil {
		panic(err)
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))
}

func httpServer() {

	r := mux.NewRouter()

	UploadRouter := r.PathPrefix("/upload").Subrouter()
	UploadRouter.HandleFunc(("/file/{shahash}/{filename}"), HashHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:5050",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
