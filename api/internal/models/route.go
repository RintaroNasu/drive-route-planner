package models

type RouteRequest struct {
	Places []string `json:"places"`
}

type Point struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

type RouteResponse struct {
	Points []Point `json:"points"`
	Route  []Point `json:"route"`
}
