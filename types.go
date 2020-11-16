package go_huawei

// RouteService is for specifying travel mode.
type RouteService string

// Avoid is for specifying routes that avoid certain features.
type Avoid int

// TransitMode is for specifying a transit mode for a request
type TransitMode string

// Travel mode preferences.
const (
	RouteServiceDriving   = RouteService("driving")
	RouteServiceWalking   = RouteService("walking")
	RouteServiceBicycling = RouteService("bicycling")
)

// Features to avoid.
const (
	AvoidTolls    = Avoid(1)
	AvoidHighways = Avoid(2)
)

// Distance is the API representation for a distance between two points.
type Distance struct {
	// HumanReadable is the human friendly distance. This is rounded and in an
	// appropriate unit for the request. The units can be overridden with a request
	// parameter.
	HumanReadable string `json:"text"`
	// Meters is the numeric distance, always in meters. This is intended to be used
	// only in algorithmic situations, e.g. sorting results by some user specified
	// metric.
	Meters int `json:"value"`
}

// TrafficMode specifies traffic prediction model when requesting future directions.
type TrafficMode int

// Traffic prediction model when requesting future directions.
const (
	TrafficModeBestGuess   = TrafficMode(0)
	TrafficModeOptimistic  = TrafficMode(2)
	TrafficModePessimistic = TrafficMode(1)
)
