package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexflint/go-arg"
	"github.com/herrewig/tomedome/go/internal/api"
	"github.com/herrewig/tomedome/go/internal/dota"
	"github.com/herrewig/tomedome/go/internal/dota/backends"
	"github.com/herrewig/tomedome/go/internal/logging"
	"github.com/sirupsen/logrus"
)

// Returns a DotaServiceApi based on the CLI args
// Possible backends include:
// - json: loads db from a json file
// - stratz: loads db from the live public Stratz GraphQL api
// - embedded: loads db from a static json file embedded in the binary
func getServiceFromArgs(log *logrus.Entry, args Args) *dota.DotaServiceApi {
	var client dota.DotaClient

	switch args.Backend {
	case "stratz":
		if args.StratzApiKey == "" {
			log.Fatal("TOMEDOME_STRATZ_API_KEY env is required")
		}
		client = backends.NewStratzClient(log, args.StratzApiKey)
	case "json":
		if args.JsonFilePath == "" {
			log.Fatal("TOMEDOME_DB_FILEPATH env is required")
		}
		client = backends.NewJsonFileClient(log, args.JsonFilePath)
	case "embedded":
		client = backends.NewEmbeddedDataClient(log, args.EmbeddedFileName)
	default:
		log.Fatalf("unknown backend: %s", args.Backend)
	}
	return dota.NewDotaService(log, client)
}

// CLI args
type Args struct {
	// Specify which data backend to use
	Backend string `arg:"--backend,required,help:--backend=[json|stratz|embedded]"`
	// Stratz.com API key
	StratzApiKey string `arg:"--stratz-api-key,env:TOMEDOME_STRATZ_API_KEY"`
	// Path to the json db file
	JsonFilePath string `arg:"--json-db-filepath,env:TOMEDOME_DB_FILEPATH"`
	LogLevel     string `arg:"env:LOGLEVEL" default:"info"`
	// For the embedded backend, the name of the embedded json file
	EmbeddedFileName string `arg:"--embedded-filepname,env:TOMEDOME_EMBEDDED_FILENAME" default:"mock_data.json"`
	// Run the API server on :8080
	RunServer bool `arg:"--run-server" help:"Run the API server"`
	// Dump the hero data as JSON to stdout
	Dump bool `arg:"--dump" help:"Dump the hero data as JSON"`
	// Enable rate limiting
	RateLimiting bool `arg:"--rate-limiting,env:TOMEDOME_RATELIMITING" default:"true" help:"Enable rate limiting"`
	// Run app in local dev mode (human readable logs vs json)
	LocalDev bool `arg:"--local-dev,env:LOCALDEV" help:"Run the server in local dev mode"`
	// Dump mock data JSON to stdout
	DumpMockData bool `arg:"--dump-mock-data" help:"Dump the mock data JSON to stdout"`
}

// When we change the data schema, we need to update the mock data so we can develop against it.
// This function dumps a small subset of the hero data to stdout. Usually it's used to write
// to mock_data.json in the assets package.
func dumpMockData(log *logrus.Entry, dotes *dota.DotaServiceApi) {
	log.Info("dumping mock data")
	heroes, err := dotes.GetAllHeroes()
	if err != nil {
		log.WithField("error", err).Fatal("failed to dump heroes")
	}

	// Shuffle the slice
	rand.Shuffle(len(heroes), func(i, j int) {
		heroes[i], heroes[j] = heroes[j], heroes[i]
	})
	// Select the first three elements
	heroesSubset := heroes[:3]

	// Serialize the subset to JSON and print to stdout
	subset, err := json.Marshal(heroesSubset)
	if err != nil {
		log.WithField("error", err).Fatal("failed to marshal subset of heroes")
	}
	fmt.Println(string(subset))
}

// Let's goooo
func main() {
	args := Args{}
	arg.MustParse(&args)
	log := logging.NewLogger(args.LogLevel, args.LocalDev)

	// Get the DotaServiceApi based on the CLI args
	dotes := getServiceFromArgs(log, args)

	// Dump a small subset of three heroes to stdout to use as mock data
	if args.DumpMockData {
		dumpMockData(log, dotes)
		return
	}

	// Dump the whole database to stdout
	if args.Dump {
		heroes, err := dotes.SerializeDb()
		if err != nil {
			log.WithField("error", err).Fatal("failed to dump heroes")
		}
		fmt.Println(string(heroes))
		return
	}

	// No action specified on CLI
	if !args.RunServer {
		log.Fatal("no action specified. use --run-server or --dump")
	}

	// Ok, based on args we must be running the api server!

	// Listen for SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a channel to listen for shutdown signals SIGINT and SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Listen for shutdown signal
	go func() {
		<-stop
		log.Info("shutdown signal received. canceling context")
		cancel()
	}()

	// Let's gooooo
	api.RunServer(ctx, log, args.RateLimiting, ":8080", dotes)
}
