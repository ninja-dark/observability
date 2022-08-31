package m

type Article struct {
	Id     string `json:"id,omitempty"`
	Title  string `json:"title"`
	UserId string `json:"userId"`
}
