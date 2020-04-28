package boot

import (
	"log"
	"net/http"

	"go-tutorial-2020/internal/config" 
	firebaseclient "go-tutorial-2020/pkg/firebaseClient"
	"go-tutorial-2020/pkg/kafka"

	"github.com/jmoiron/sqlx"

	userData "go-tutorial-2020/internal/data/user"
	server "go-tutorial-2020/internal/delivery/http"
	userHandler "go-tutorial-2020/internal/delivery/http/user"
	userService "go-tutorial-2020/internal/service/user"

	kConsumer "go-tutorial-2020/internal/delivery/kafka"
	httpclient "go-tutorial-2020/pkg/httpClient"
)

// HTTP will load configuration, do dependency injection and then start the HTTP server
func HTTP() error {
	var (
		s     server.Server        // HTTP Server Object
		ud    userData.Data        // User domain data layer
		us    userService.Service  // User domain service layer
		uh    *userHandler.Handler // User domain handler
		cfg   *config.Config       // Configuration object
		fb    *firebaseclient.Client
		k     *kafka.Kafka // Kafka Producer
		httpc *httpclient.Client
	)

	// Get configuration
	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	cfg = config.Get()

	httpc = httpclient.NewClient()

	// Open MySQL DB Connection
	db, err := sqlx.Open("mysql", cfg.Database.Master)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	fb, err = firebaseclient.NewClient(cfg)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	k, err = kafka.New(cfg.Kafka.Username, cfg.Kafka.Password, cfg.Kafka.Brokers)
	if err != nil {  
		log.Fatalf("[KAFKA] Failed to initialize kafka producer: %v", err)
	}

	// User domain initialization
	ud = userData.New(db, fb, httpc)
	us = userService.New(ud, k)
	uh = userHandler.New(us)

	// Inject service used on handler
	s = server.Server{
		User: uh,
	}


	go kConsumer.New(us, k, cfg.Kafka.Subscriptions)
	
	// Error Handling
	if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		return err
	}

	return nil
}
