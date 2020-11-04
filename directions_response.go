package go_huawei

type DirectionsResponse struct {
	Routes []Route `json:"routes"`
	CommonResponse
}

type Route struct {
	Paths  []Path           `json:"paths"`
	Bounds CoordinateBounds `json:"bounds"`
}

type CoordinateBounds struct {
	Southwest Coordinate `json:"southwest"`
	Northeast Coordinate `json:"northeast"`
}

type Path struct {
	Duration              float64    `json:"duration"`
	DurationText          string     `json:"durationText"`
	DurationInTrafficText string     `json:"durationInTrafficText"`
	DurationInTraffic     float64    `json:"durationInTraffic"`
	Distance              float64    `json:"distance"`
	StartLocation         Coordinate `json:"startLocation"`
	StartAddress          string     `json:"startAddress"`
	DistanceText          string     `json:"distanceText"`
	Steps                 []Step     `json:"steps"`
	EndLocation           Coordinate `json:"endLocation"`
	EndAddress            string     `json:"endAddress"`
}

type Step struct {
	Duration      float64      `json:"duration"`
	Orientation   int64        `json:"orientation"`
	DurationText  string       `json:"durationText"`
	Distance      float64      `json:"distance"`
	StartLocation Coordinate   `json:"startLocation"`
	Instruction   string       `json:"instruction"`
	Action        Action       `json:"action"`
	DistanceText  string       `json:"distanceText"`
	EndLocation   Coordinate   `json:"endLocation"`
	Polyline      []Coordinate `json:"polyline"`
	RoadName      string       `json:"roadName"`
}

type Action string

const (
	End             Action = "end"
	ForkLeft        Action = "fork-left"
	ForkRight       Action = "fork-right"
	RampRight       Action = "ramp-right"
	RampLeft        Action = "ramp-left"
	RoundaboutRight Action = "roundabout-right"
	RoundaboutLeft  Action = "roundabout-left"
	Straight        Action = "straight"
	TurnLeft        Action = "turn-left"
	TurnRight       Action = "turn-right"
	TurnSlightLeft  Action = "turn-slight-left"
	TurnSlightRight Action = "turn-slight-right"
)
