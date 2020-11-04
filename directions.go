package go_huawei

import (
	"context"
	"errors"
	"fmt"
)

var directionsAPI = &apiConfig{
	host:          "https://mapapi.cloud.huawei.com",
	path:          "/mapApi/v1/routeService/",
	acceptsApiKey: true,
}

// Directions issues the Directions request and retrieves the Response
func (c *Client) Directions(ctx context.Context, directionsRequest *DirectionsRequest) ([]Route, error) {
	if directionsRequest.Origin.isEmpty() {
		return nil, errors.New("map-kit: origin missing")
	}

	if directionsRequest.Destination.isEmpty() {
		return nil, errors.New("map-kit: destination missing")
	}
	if directionsRequest.RouteService != "" &&
		RouteServiceDriving != directionsRequest.RouteService &&
		RouteServiceWalking != directionsRequest.RouteService &&
		RouteServiceBicycling != directionsRequest.RouteService {
		return nil, fmt.Errorf("map-kit: unknown RouteService: '%s'", directionsRequest.RouteService)
	}

	response := DirectionsResponse{}

	if err := c.postJSON(ctx, directionsAPI, directionsRequest, &response, directionsRequest.RouteService); err != nil {
		return nil, err
	}

	if err := response.StatusError(); err != nil {
		return nil, err
	}

	return response.Routes, nil
}
