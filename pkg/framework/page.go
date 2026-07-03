package framework

import "net/http"

type Page struct {
	Status      int
	Title       string
	Description string
	Template    string
	Data        any
}

func (p Page) render(w http.ResponseWriter) {

}
