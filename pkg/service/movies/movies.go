package movies

type MovieReq struct {
	Title      string `json:"title"`
	Genre      string `json:"genre"`
	ReleasDate string `json:"release_date"`
	Duration   string `json:"duration"`
}

func CreateMovies() {

}
