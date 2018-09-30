package controller

import (
	"fmt"
	"net/http"
)

type HomeController struct {
}

func (hc HomeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Wellcome to air traffic control")
}
