package go_huawei

import (
	"fmt"
	"strconv"
	"strings"
)

type Coordinate struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

func (c *Coordinate) String() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("%f,%f", c.Lat, c.Lng)
}

func (c *Coordinate) isEmpty() bool {
	return c == nil
}

func ParseCoordinate(str string) *Coordinate {
	elements := strings.Split(str, ",")
	if len(elements) != 2 {
		return nil
	}

	lat, err := strconv.ParseFloat(elements[0], 64)
	if err != nil {
		return nil
	}

	lng, err := strconv.ParseFloat(elements[1], 64)
	if err != nil {
		return nil
	}

	return &Coordinate{
		Lng: lng,
		Lat: lat,
	}
}
