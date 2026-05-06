package service

import (
	"errors"
	"testing"

	"github.com/RintaroNasu/drive-route-planner/api/internal/httpx"
	"github.com/RintaroNasu/drive-route-planner/api/internal/models"
	"github.com/stretchr/testify/require"
)

func TestRouteService_Normal(t *testing.T) {

	tests := []struct {
		name  string
		input []string
		mock  map[string]models.Point
	}{
		{
			name:  "【正常系】1地点のみ",
			input: []string{"東京タワー"},
			mock: map[string]models.Point{
				"東京タワー": {Name: "東京タワー", Lat: 1, Lng: 1},
			},
		},
		{
			name:  "【正常系】2地点",
			input: []string{"A", "B"},
			mock: map[string]models.Point{
				"A": {Name: "A", Lat: 0, Lng: 0},
				"B": {Name: "B", Lat: 1, Lng: 1},
			},
		},
		{
			name:  "【正常系】複数地点（3以上）",
			input: []string{"A", "B", "C"},
			mock: map[string]models.Point{
				"A": {Name: "A", Lat: 0, Lng: 0},
				"B": {Name: "B", Lat: 5, Lng: 5},
				"C": {Name: "C", Lat: 1, Lng: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := newRouteService(func(place string) (models.Point, error) {
				return tt.mock[place], nil
			})

			res, err := svc.RouteFromPlace(tt.input)

			require.NoError(t, err)
			require.NotNil(t, res)

			require.NotNil(t, res.Points)
			require.NotNil(t, res.Route)
			require.IsType(t, []models.Point{}, res.Points)
			require.IsType(t, []models.Point{}, res.Route)

			require.Len(t, res.Points, len(tt.input))
			require.Len(t, res.Route, len(tt.input))

			if len(tt.input) == 1 {
				require.Equal(t, res.Points, res.Route)
			}

			for _, p := range res.Route {
				require.NotEmpty(t, p.Name)
			}
		})
	}
}

func TestRouteService_Error(t *testing.T) {

	tests := []struct {
		name     string
		input    []string
		mockFunc func(string) (models.Point, error)
		wantCode string
	}{
		{
			name:     "【異常系】places空",
			input:    []string{},
			wantCode: "InvalidRequest",
		},
		{
			name:  "【異常系】地名が見つからない",
			input: []string{"unknown"},
			mockFunc: func(place string) (models.Point, error) {
				return models.Point{}, httpx.NotFound("place not found", nil)
			},
			wantCode: "NotFound",
		},
		{
			name:  "【異常系】外部APIエラー",
			input: []string{"A"},
			mockFunc: func(place string) (models.Point, error) {
				return models.Point{}, httpx.ExternalAPI("request failed", errors.New("timeout"))
			},
			wantCode: "ExternalAPIError",
		},
		{
			name:  "【異常系】JSONパースエラー",
			input: []string{"A"},
			mockFunc: func(place string) (models.Point, error) {
				return models.Point{}, httpx.Internal("json parse failed", errors.New("decode error"))
			},
			wantCode: "InternalError",
		},
		{
			name:  "【異常系】lat/lngパースエラー",
			input: []string{"A"},
			mockFunc: func(place string) (models.Point, error) {
				return models.Point{}, httpx.Internal("lat parse failed", errors.New("parse error"))
			},
			wantCode: "InternalError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := newRouteService(tt.mockFunc)

			res, err := svc.RouteFromPlace(tt.input)

			require.Error(t, err)
			require.Nil(t, res)

			appErr, ok := err.(*httpx.AppError)
			require.True(t, ok)

			require.Equal(t, tt.wantCode, appErr.Code)
		})
	}
}

func TestBuildRoute(t *testing.T) {

	tests := []struct {
		name  string
		input []models.Point
		want  []string
	}{
		{
			name: "【アルゴリズム】順序が変わる",
			input: []models.Point{
				{Name: "A", Lat: 0, Lng: 0},
				{Name: "B", Lat: 5, Lng: 5},
				{Name: "C", Lat: 1, Lng: 1},
			},
			want: []string{"A", "C", "B"},
		},
		{
			name: "【アルゴリズム】同距離（順序維持）",
			input: []models.Point{
				{Name: "A", Lat: 0, Lng: 0},
				{Name: "B", Lat: 1, Lng: 1},
				{Name: "C", Lat: 1, Lng: 1},
			},
			want: []string{"A", "B", "C"},
		},
		{
			name: "【アルゴリズム】1件のみ",
			input: []models.Point{
				{Name: "A", Lat: 0, Lng: 0},
			},
			want: []string{"A"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res := buildRoute(tt.input)

			require.Len(t, res, len(tt.want))

			for i, p := range res {
				require.Equal(t, tt.want[i], p.Name)
			}
		})
	}
}
