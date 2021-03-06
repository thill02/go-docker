package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	// "api/root/appconfig"
	"github.com/thill02/go-docker/dbutils"
	"github.com/thill02/go-docker/handlers"

	"github.com/apex/gateway"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "github.com/lib/pq"
)

func lambdaContext() bool {

	value, exists := os.LookupEnv("LAMBDA")

	if exists && strings.ToUpper(value) == "TRUE" {
		return true
	}
	return false
}

func main() {
	//  Here's what would typically be here:
	// lambda.Start(Handler)
	//
	// There were a few options on how to incorporate Echo v4 on Lambda.
	//
	// Landed here for now:
	//
	//     https://github.com/apex/gateway
	//     https://github.com/labstack/echo/issues/1195
	//
	// With this for local development:
	//     https://medium.com/a-man-with-no-server/running-go-aws-lambdas-locally-with-sls-framework-and-sam-af3d648d49cb
	//
	// This looks promising and is from awslabs, but Echo v4 support isn't quite there yet.
	// There is a pull request in progress, Re-evaluate in April 2020.
	//
	//    https://github.com/awslabs/aws-lambda-go-api-proxy
	//

	db := dbutils.Connection()

	e := echo.New()
	e.Use(middleware.CORS())
	// e.Use(middleware.JWTWithConfig(appconfig.JWTConfig))

	// Routes
	// Instrument Groups
	e.GET("instrument_groups", handlers.ListInstrumentGroups(db))
	e.POST("instrument_groups", handlers.CreateInstrumentGroup(db))
	e.GET("instrument_groups/:id", handlers.GetInstrumentGroup(db))
	e.DELETE("instrument_groups/:id", handlers.DeleteInstrumentGroup(db))
	// Instrument Group Instruments
	e.GET("instrument_groups/:id/instruments", handlers.ListInstrumentGroupInstruments(db))
	// Add instrument to group
	e.POST("instrument_groups/:id/instruments", handlers.CreateInstrumentGroupInstruments(db))
	// Remove instrument from group
	e.DELETE("instrument_groups/:id/instruments/:instrument_id", handlers.DeleteInstrumentGroupInstruments(db))
	// Instruments
	e.GET("instruments", handlers.ListInstruments(db))
	e.POST("instruments", handlers.CreateInstrument(db))
	e.GET("instruments/:id", handlers.GetInstrument(db))
	// e.GET("instruments/:id/timeseries", handlers.GetInstrument(db))

	e.PUT("instruments/:id", handlers.UpdateInstrument(db))
	e.DELETE("instruments/:id", handlers.DeleteInstrument(db))

	// Time Series
	e.GET("timeseries", handlers.GetTimeseries(db))
	e.GET("timeseries/:id", handlers.GetTimeseries(db))
	// e.GET("timeseries/:id/measurements", handlers.GetTimeseriesMeasurements(db))
	// ^ Query Parameters
	e.GET("timeseries_measurements", handlers.GetTimeseriesMeasurements(db))
	// Domains
	e.GET("domains", handlers.GetDomains(db))

	log.Printf(
		"starting server; Running On AWS LAMBDA: %t",
		lambdaContext(),
	)

	if lambdaContext() {
		log.Fatal(gateway.ListenAndServe(":3030", e))
	} else {
		log.Fatal(http.ListenAndServe(":3030", e))
	}
}