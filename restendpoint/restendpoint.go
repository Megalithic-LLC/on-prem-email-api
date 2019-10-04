package restendpoint

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/on-prem-net/email-api/agentstreamendpoint"
	"github.com/docktermj/go-logger/logger"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type RestEndpoint struct {
	agentStreamEndpoint *agentstreamendpoint.AgentStreamEndpoint
	authMiddleware      AuthenticationMiddleware
	db                  *gorm.DB
	listener            net.Listener
	redisClient         *redis.Client
	router              *mux.Router
	server              *http.Server
}

func New(
	agentStreamEndpoint *agentstreamendpoint.AgentStreamEndpoint,
	authMiddleware *AuthenticationMiddleware,
	db *gorm.DB,
	redisClient *redis.Client,
) *RestEndpoint {

	router := mux.NewRouter()

	self := RestEndpoint{
		agentStreamEndpoint: agentStreamEndpoint,
		authMiddleware:      *authMiddleware,
		db:                  db,
		redisClient:         redisClient,
		router:              router,
	}

	router.HandleFunc("/v1/agentStream", self.agentStream).Methods("GET")

	router.HandleFunc("/v1/agents", self.getAgents).Methods("GET")
	router.HandleFunc("/v1/agents", self.createAgent).Methods("POST")
	router.HandleFunc("/v1/agents/{id}", self.getAgent).Methods("GET")

	router.HandleFunc("/v1/emailAccounts", self.getEmailAccounts).Methods("GET")
	router.HandleFunc("/v1/emailAccounts", self.createEmailAccount).Methods("POST")
	router.HandleFunc("/v1/emailAccounts/{id}", self.deleteEmailAccount).Methods("DELETE")
	router.HandleFunc("/v1/emailAccounts/{id}", self.getEmailAccount).Methods("GET")

	router.HandleFunc("/v1/plans", self.getPlans).Methods("GET")
	router.HandleFunc("/v1/plans/{id}", self.getPlan).Methods("GET")

	router.HandleFunc("/v1/services", self.getServices).Methods("GET")
	router.HandleFunc("/v1/services/{id}", self.getService).Methods("GET")

	router.HandleFunc("/v1/serviceInstances", self.getServiceInstances).Methods("GET")
	router.HandleFunc("/v1/serviceInstances", self.createServiceInstance).Methods("POST")
	router.HandleFunc("/v1/serviceInstances/{id}", self.getServiceInstance).Methods("GET")

	router.HandleFunc("/v1/snapshots", self.getSnapshots).Methods("GET")
	router.HandleFunc("/v1/snapshots", self.createSnapshot).Methods("POST")
	router.HandleFunc("/v1/snapshots/{id}", self.getSnapshot).Methods("GET")

	router.HandleFunc("/v1/tokenAuth", self.createToken).Methods("POST")
	router.HandleFunc("/v1/tokenRefresh", self.refreshToken).Methods("POST")

	router.Use(self.authMiddleware.Middleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.Handle("/", router)
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		logger.Fatalf("Failed listening: %v", err)
	}
	self.server = &http.Server{}

	go func() {
		if err := self.server.Serve(listener); err != nil {
			logger.Errorf("Failed listening: %v", err)
		}
	}()

	logger.Infof("Listening for http on port %d", listener.Addr().(*net.TCPAddr).Port)
	return &self
}

func (self *RestEndpoint) Shutdown(ctx context.Context) {
	self.server.Shutdown(ctx)
	logger.Infof("Rest endpoint shutdown")
}
