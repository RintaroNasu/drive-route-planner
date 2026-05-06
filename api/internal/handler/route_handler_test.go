package handler

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RintaroNasu/drive-route-planner/api/internal/httpx"
	"github.com/RintaroNasu/drive-route-planner/api/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mockRouteService struct {
	RouteFromPlaceFunc func([]string) (*models.RouteResponse, error)
	called             int
}

func (m *mockRouteService) RouteFromPlace(places []string) (*models.RouteResponse, error) {
	m.called++
	return m.RouteFromPlaceFunc(places)
}
func TestRouteHandler(t *testing.T) {

	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	tests := []struct {
		name           string
		reqBody        string
		svc            mockRouteService
		wantStatusCode int
		wantBodyPart   string
		wantCalled     int
	}{
		{
			name:    "【正常系】正しいリクエストなら200が返る",
			reqBody: `{"places":["A"]}`,
			svc: mockRouteService{
				RouteFromPlaceFunc: func(places []string) (*models.RouteResponse, error) {
					require.Equal(t, []string{"A"}, places)
					return &models.RouteResponse{
						Route: []models.Point{
							{Name: "A", Lat: 1, Lng: 1},
						},
					}, nil
				},
			},
			wantStatusCode: http.StatusOK,
			wantBodyPart:   `"name":"A"`,
			wantCalled:     1,
		},
		{
			name:    "【異常系】serviceエラーならそのまま500",
			reqBody: `{"places":["A"]}`,
			svc: mockRouteService{
				RouteFromPlaceFunc: func(places []string) (*models.RouteResponse, error) {
					return nil, httpx.Internal("error", errors.New("fail"))
				},
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBodyPart:   `"code":"InternalError"`,
			wantCalled:     1,
		},
		{
			name:    "【異常系】不正JSONなら400（serviceは呼ばれない）",
			reqBody: `{"places":}`,
			svc: mockRouteService{
				RouteFromPlaceFunc: func(places []string) (*models.RouteResponse, error) {
					return nil, nil
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantBodyPart:   `"code":"InvalidRequest"`,
			wantCalled:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := tt.svc
			h := NewRouteHandler(&svc)

			req := httptest.NewRequest(http.MethodPost, "/route-from-place",
				bytes.NewBuffer([]byte(tt.reqBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.RouteFromPlace(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatusCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBodyPart)
			require.Equal(t, tt.wantCalled, svc.called)
		})
	}
}
