// +build !enterprise

package api

import (
	"log"
	"net/http"
	"os"
)

var _ = setupPayment()

func setupPayment() bool {
	log.Println("PAYMENT IS ACTIVE!")
	os.Setenv("IS_ENTERPRISE", "false")

	http.HandleFunc("/payment", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("There is payment!"))
	})
	return false
}
