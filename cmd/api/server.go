package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ankeshnirala/order-mgmt/auth-service/internal/routers"
	"github.com/ankeshnirala/order-mgmt/common-service/middleware"
	"github.com/ankeshnirala/order-mgmt/common-service/storage"
	"github.com/ankeshnirala/order-mgmt/common-service/utils/shutdown"
)

type Server struct {
	logger     *logrus.Logger
	listenAddr string
	sqlStore   storage.MySQLStorage
}

func NewServer(logger *logrus.Logger, listenAddr string, sqlStore storage.MySQLStorage) *Server {
	return &Server{
		logger:     logger,
		listenAddr: listenAddr,
		sqlStore:   sqlStore,
	}
}

func (s *Server) Start() error {
	// initialize new router
	router := mux.NewRouter().StrictSlash(true)
	midd := middleware.New(s.logger, s.sqlStore)
	router.Use(midd.RequestLogger)

	routers.New(s.logger, router, s.sqlStore).RegisterAuthRoutes()

	// added MethodNotAllowedHandler & PathNotFoundHanler middleware
	routers.New(s.logger, router, s.sqlStore).RegisterDefaultMiddleware()

	// create http server instance
	srv := &http.Server{
		Addr:         s.listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			msg := fmt.Sprintf("error - server.go: %v", err)
			s.logger.Error(msg)
			return
		}
	}()

	// ensure the server is shutdown gracefully & app runs
	shutdown.Gracefully(s.logger, srv)

	return nil
}
