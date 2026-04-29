package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RintaroNasu/drive-route-planner/api/internal/models"
)

type RouteService struct{}

func NewRouteService() *RouteService {
	return &RouteService{}
}

func (s *RouteService) RouteFromPlace(places []string) (*models.RouteResponse, error) {
	if len(places) == 0 {
		return nil, errors.New("places is required")
	}

	var points []models.Point

	for _, place := range places {
		p, err := s.geocode(place)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}

	return &models.RouteResponse{
		Points: points,
		Route:  points,
	}, nil
}

func (s *RouteService) geocode(place string) (models.Point, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"

	params := url.Values{}
	params.Add("q", place)
	params.Add("format", "json")
	params.Add("limit", "1")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "go-route-app")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Point{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return models.Point{}, fmt.Errorf("external api error: %d", resp.StatusCode)
	}

	var result []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.Point{}, err
	}

	if len(result) == 0 {
		return models.Point{}, errors.New("place not found")
	}

	var lat, lng float64
	fmt.Sscanf(result[0].Lat, "%f", &lat)
	fmt.Sscanf(result[0].Lon, "%f", &lng)

	return models.Point{
		Name: place,
		Lat:  lat,
		Lng:  lng,
	}, nil
}
