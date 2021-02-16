package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	userapi "github.com/sm43/goa-gorm"
	user "github.com/sm43/goa-gorm/gen/user"
	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	// Blank for package side effect: loads postgres drivers
	_ "github.com/lib/pq"
)

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[userapi] ", log.Ltime)
	}

	// Database Connection
	var (
		db *gorm.DB
	)
	{
		if err := godotenv.Load(".env.dev"); err != nil {
			fmt.Fprintf(os.Stderr, "SKIP: loading env file %s failed: %s\n", "file", err)
		}
		viper.AutomaticEnv()

		dbConfig, err := initDB()
		if err != nil {
			log.Fatal("failed to get db configurations", err)
		}

		conn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

		db, err = gorm.Open("postgres", conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successful Db Connection")
		defer db.Close()

		db.AutoMigrate(userapi.User{})
	}
	// Initialize the services.
	var (
		userSvc user.Service
	)
	{
		userSvc = userapi.NewUser(db, logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		userEndpoints *user.Endpoints
	)
	{
		userEndpoints = user.NewEndpoints(userSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "http://localhost:8080"
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h := strings.Split(u.Host, ":")[0]
				u.Host = h + ":" + *httpPortF
			} else if u.Port() == "" {
				u.Host += ":80"
			}
			handleHTTPServer(ctx, u, userEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		fmt.Fprintf(os.Stderr, "invalid host argument: %q (valid hosts: localhost)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}

func initDB() (*Database, error) {

	db := &Database{}
	if db.Host = viper.GetString("POSTGRES_HOST"); db.Host == "" {
		return nil, fmt.Errorf("no POSTGRES_HOST environment variable defined")
	}
	if db.Port = viper.GetString("POSTGRES_PORT"); db.Port == "" {
		return nil, fmt.Errorf("no POSTGRES_PORT environment variable defined")
	}
	if db.Name = viper.GetString("POSTGRES_DB"); db.Name == "" {
		return nil, fmt.Errorf("no POSTGRES_DB environment variable defined")
	}
	if db.User = viper.GetString("POSTGRES_USER"); db.User == "" {
		return nil, fmt.Errorf("no POSTGRES_USER environment variable defined")
	}
	if db.Password = viper.GetString("POSTGRES_PASSWORD"); db.Password == "" {
		return nil, fmt.Errorf("no POSTGRES_PASSWORD environment variable defined")
	}
	return db, nil
}
