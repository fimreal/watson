package serve

import (
	"fmt"
	"net/http"

	"github.com/fimreal/goutils/ezap"
)

func Run(port string) {
	http.HandleFunc("/hc",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})

	ezap.Info("listening at ", port)
	ezap.Fatal(http.ListenAndServe(port, nil))
}
