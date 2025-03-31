package models

// Маршрут: /api/address/search метод POST
type SearchRequest struct {
	Query string `json:"query"`
}

// Маршрут: /api/address/geocode метод POST
type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// Маршрут: /api/register метод POST
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Маршрут: /api/login метод POST
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
