package go_huawei

import "encoding/json"

func (r *DirectionsRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DirectionsRequest struct {
	Origin      *Coordinate   `json:"origin"`
	Destination *Coordinate   `json:"destination"`
	Waypoints   []*Coordinate `json:"waypoints,omitempty"`

	//Indicates whether to optimize the waypoint. false (default): no true: yes
	Optimize bool `json:"optimize,omitempty"`

	//Time estimation mode. The options are as follows: 0 (default): best guess 1: The traffic condition is worse than the
	//historical average. 2: The traffic condition is better than the historical average
	TrafficMode TrafficMode `json:"trafficMode,omitempty"`

	// Waypoint type. The options are as follows: false (default): stopover true: via (pass-by)
	ViaType bool `json:"viaType,omitempty"`

	// Indicates whether to return multiple planned routes. The options are as follows: true: yes false (default): no Note:
	//This parameter is unavailable when waypoints are set.
	Alternatives bool `json:"alternatives,omitempty"`

	//Indicates the specified type of roads to be avoided. The options are as follows: 1: Avoid toll roads. 2: Avoid
	//expressways. If this parameter is not included in the request, the route taking the least time will be returned by
	//default
	Avoid []Avoid `json:"avoid,omitempty"`

	//Estimated departure time, in seconds since 00:00 on January 1, 1970 (UTC). The value must be the current time or a
	//future time.
	DepartAt int `json:"departAt,omitempty"`

	// Language of the distance and journey time descriptions in the returned result. Currently, only zh_CN and en are supported
	Language string `json:"language,omitempty"`

	RouteService RouteService `json:"-"`
}
