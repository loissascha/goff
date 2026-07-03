package framework

type Page struct {
	Status      int
	Title       string
	Description string
	Template    string
	Data        any
}
