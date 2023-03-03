package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"

	"github.com/stremovskyy/go-huawei"
)

var (
	apiKey       = flag.String("key", "", "API Key for using Huawei Map-kit API.")
	origin       = flag.String("origin", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
	destination  = flag.String("destination", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
	waypoints    = flag.String("waypoints", "", "The waypoints for driving directions request, | separated.")
	alternatives = flag.Bool("alternatives", false, "Whether the Directions service may provide more than one route alternative in the response.")
	trafficModel = flag.String("traffic_model", "", "Specifies traffic prediction model when request future directions. Valid values are optimistic, best_guess, and pessimistic. Optional.")
	avoid        = flag.String("avoid", "", "Indicates that the calculated route(s) should avoid the indicated features, | separated.")
	language     = flag.String("language", "", "Specifies the language in which to return results.")
	routeService = flag.String("routeService", "driving", "The travel routeService for this directions request.")
)

func usageAndExit(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func main() {
	flag.Parse()

	var client *go_huawei.Client
	var err error
	if *apiKey != "" {
		client, err = go_huawei.NewClient(
			go_huawei.WithAPIKey(*apiKey),
			go_huawei.WithRateLimit(2),
		)
	} else {
		usageAndExit("Please specify an API Key.")
	}
	check(err)

	request := &go_huawei.DirectionsRequest{
		Origin:       go_huawei.ParseCoordinate(*origin),
		Destination:  go_huawei.ParseCoordinate(*destination),
		Alternatives: *alternatives,
		Language:     *language,
	}

	lookupRouteService(*routeService, request)
	lookupTrafficModel(*trafficModel, request)

	if *waypoints != "" {
		stringPoints := strings.Split(*waypoints, "|")
		for _, stringPoint := range stringPoints {
			request.Waypoints = append(request.Waypoints, go_huawei.ParseCoordinate(stringPoint))
		}
	}

	if *avoid != "" {
		for _, a := range strings.Split(*avoid, "|") {
			switch a {
			case "tolls":
				request.Avoid = append(request.Avoid, go_huawei.AvoidTolls)
			case "highways":
				request.Avoid = append(request.Avoid, go_huawei.AvoidHighways)
			default:
				log.Fatalf("Unknown avoid restriction %s", a)
			}
		}
	}

	routes, err := client.Directions(context.Background(), request)
	check(err)

	_, _ = pretty.Println(waypoints)

	for i, route := range routes {
		for k, path := range route.Paths {
			fmt.Printf("%d:%d-> %s\n%s\n", i, k, path.DurationInTrafficText, string(path.OverviewPolyline()))
		}
	}
}

func lookupRouteService(routeService string, r *go_huawei.DirectionsRequest) {
	switch routeService {
	case "driving":
		r.RouteService = go_huawei.RouteServiceDriving
	case "walking":
		r.RouteService = go_huawei.RouteServiceWalking
	case "bicycling":
		r.RouteService = go_huawei.RouteServiceBicycling
	case "":
		// ignore
	default:
		log.Fatalf("Unknown routeService '%s'", routeService)
	}
}

func lookupTrafficModel(trafficModel string, r *go_huawei.DirectionsRequest) {
	switch trafficModel {
	case "optimistic":
		r.TrafficMode = go_huawei.TrafficModeOptimistic
	case "best_guess":
		r.TrafficMode = go_huawei.TrafficModeBestGuess
	case "pessimistic":
		r.TrafficMode = go_huawei.TrafficModePessimistic
	case "":
		// ignore
	default:
		log.Fatalf("Unknown traffic routeService %s", trafficModel)
	}
}
