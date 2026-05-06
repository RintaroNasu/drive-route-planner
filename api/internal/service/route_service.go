package service

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/RintaroNasu/drive-route-planner/api/internal/httpx"
	"github.com/RintaroNasu/drive-route-planner/api/internal/models"
)

type RouteService interface {
	RouteFromPlace(places []string) (*models.RouteResponse, error)
}

type routeService struct {
	geocodeFunc func(string) (models.Point, error)
}

func NewRouteService() RouteService {
	return newRouteService(nil)
}

func newRouteService(geocodeFn func(string) (models.Point, error)) *routeService {
	s := &routeService{}
	if geocodeFn != nil {
		s.geocodeFunc = geocodeFn
	} else {
		s.geocodeFunc = s.geocode
	}
	return s
}

func (s *routeService) RouteFromPlace(places []string) (*models.RouteResponse, error) {
	if len(places) == 0 {
		return nil, httpx.InvalidRequest("places is required", nil)
	}

	var points []models.Point

	for _, place := range places {
		p, err := s.geocodeFunc(place)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}

	return &models.RouteResponse{
		Points: points,
		Route:  buildRoute(points),
	}, nil
}

func (s *routeService) geocode(place string) (models.Point, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"

	params := url.Values{}
	params.Add("q", place)
	params.Add("format", "json")
	params.Add("limit", "1")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return models.Point{}, httpx.Internal("request creation failed", err)
	}
	req.Header.Set("User-Agent", "go-route-app")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.Point{}, httpx.ExternalAPI("request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return models.Point{}, httpx.ExternalAPI(fmt.Sprintf("external api error: %d", resp.StatusCode), nil)
	}

	var result []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.Point{}, httpx.Internal("json parse failed", err)
	}

	if len(result) == 0 {
		return models.Point{}, httpx.NotFound("place not found", nil)
	}

	var lat, lng float64
	if _, err := fmt.Sscanf(result[0].Lat, "%f", &lat); err != nil {
		return models.Point{}, httpx.Internal("lat parse failed", err)
	}
	if _, err := fmt.Sscanf(result[0].Lon, "%f", &lng); err != nil {
		return models.Point{}, httpx.Internal("lng parse failed", err)
	}

	return models.Point{
		Name: place,
		Lat:  lat,
		Lng:  lng,
	}, nil
}

func buildRoute(points []models.Point) []models.Point {
	if len(points) == 0 {
		return points
	}

	visited := make([]bool, len(points))
	var route []models.Point

	// スタート地点（最初は0番目の地点）
	current := 0
	route = append(route, points[current])
	visited[current] = true

	for len(route) < len(points) {
		// 次に訪問する地点のインデックス（未決定なので-1）
		next := -1

		// 現在地点から最も近い地点を探すための最小距離（初期値は最大値）
		minDist := math.MaxFloat64

		for i := 0; i < len(points); i++ {
			if visited[i] {
				continue
			}

			d := distance(points[current], points[i])
			if d < minDist {
				minDist = d
				next = i
			}
		}

		visited[next] = true
		route = append(route, points[next])
		current = next
	}

	return route
}

func distance(a, b models.Point) float64 {
	dx := a.Lat - b.Lat
	dy := a.Lng - b.Lng

	// ユークリッド距離（平方根は不要：比較のみのため）
	return dx*dx + dy*dy
}
