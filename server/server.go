package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/jgarrone/foaas-limiting/services/foaasapi"
	log "github.com/sirupsen/logrus"
)

const (
	MessagePath  = "/message"
	UserIdHeader = "userId"
)

type Server struct {
	address        string
	limiter        Limiter
	messageService foaasapi.Service
	stopCh         chan os.Signal
}

func New(address string, limiter Limiter, messageService foaasapi.Service) *Server {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return &Server{
		address:        address,
		limiter:        limiter,
		messageService: messageService,
		stopCh:         stopCh,
	}
}

func (s *Server) Run() {
	router := mux.NewRouter()
	router.HandleFunc(MessagePath, s.handleMessage).Methods("GET")

	ctx := context.Background()
	sv := &http.Server{Addr: s.address, Handler: router}

	go func() {
		if err := sv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting the server: %v", err)
		}
	}()
	log.Infof("Started listening at %s", s.address)

	// Wait here until a termination signal is received.
	<-s.stopCh
	log.Info("Server stopped")

	// Gracefully shutdown the server.
	if err := sv.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}
}

func (s *Server) Stop() {
	s.stopCh <- syscall.SIGINT
}

func (s *Server) handleMessage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		s.respond(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userId := req.Header.Get(UserIdHeader)
	if userId == "" {
		s.respond(w, http.StatusUnauthorized, fmt.Sprintf("must provide header %s", UserIdHeader))
		return
	}

	if !utf8.ValidString(userId) {
		s.respond(w, http.StatusBadRequest, fmt.Sprintf("%s header must containg a vaild utf-8 value", UserIdHeader))
		return
	}

	if !s.limiter.AllowRequestFrom(userId) {
		s.respond(w, http.StatusTooManyRequests, "quota for user exceeded")
		return
	}

	resp, err := s.messageService.GetMessageFor(userId)
	if err != nil {
		log.Errorf("error getting message for %q, err: %v", userId, err)
		s.respond(w, http.StatusInternalServerError, "error fetching message, try again later")
		return
	}

	s.respond(w, http.StatusOK, resp.Message)
}

func (s *Server) respond(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := map[string]string{
		"message": message,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("error encoding response: %v", err)
		return
	}

	if _, err := w.Write(jsonResp); err != nil {
		log.Errorf("error writing response: %v", err)
	}
}
