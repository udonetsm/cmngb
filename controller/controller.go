package controller

import (
	"fmt"
	"io"
	"net/http"
)

func UpdateNumberController(w http.ResponseWriter, r *http.Request) {
	Request(w, r, "update number")
}

func UpdateNameController(w http.ResponseWriter, r *http.Request) {
	Request(w, r, "update name")

}

func UpdateNumListController(w http.ResponseWriter, r *http.Request) {
	Request(w, r, "update number list")
}

func InfoController(w http.ResponseWriter, r *http.Request) {
	Request(w, r, "info")
}

func DeleteController(w http.ResponseWriter, r *http.Request) {
	Request(w, r, "delete")
}

func CreateController(w http.ResponseWriter, r *http.Request) {
	Request(w, r, "create")
}

func Request(w http.ResponseWriter, r *http.Request, action string) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(string(req))
	w.Write([]byte("HELLO" + action))
}
