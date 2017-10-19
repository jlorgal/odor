package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jlorgal/odor/odor"
	"github.com/jlorgal/odor/odor/svc"
)

const (
	// ServiceName is the name of the service
	ServiceName = "odor"
	// Component is the name of the component
	Component = "profile"
)

var profiles map[string]*odor.Profile

// Service represents the Profile service
type Service struct {
	*odor.Config
	server *http.Server
	users  chan *odor.Profile
}

// New creates
func New(config *odor.Config) *Service {
	s := &Service{}
	s.Config = config
	return s
}

// Start the service
func (s *Service) Start() error {
	if s.server != nil {
		return fmt.Errorf("Service is already started")
	}
	s.server = &http.Server{
		Addr:    s.Address,
		Handler: s.router(),
	}
	s.users = make(chan *odor.Profile)
	profiles = map[string]*odor.Profile{}
	go s.run()
	return s.server.ListenAndServe()
}

// Stop the service
func (s *Service) Stop() error {
	if s.server != nil {
		close(s.users)
		s.server.Shutdown(nil)
		s.server = nil
	}
	return nil
}

func (s *Service) router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/{msisdn}", s.withMws("updateUserProfile")(s.UpdateUserProfile)).Methods("PUT")
	r.HandleFunc("/", s.withMws("Welcome")(s.Welcome)).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(s.withMws("notFound")(svc.WithNotFound()))
	return r
}

func (s *Service) withMws(op string) func(http.HandlerFunc) http.HandlerFunc {
	logContext := &svc.LogContext{
		Service:   ServiceName,
		Component: Component,
		Operation: op,
	}
	withLogContext := svc.WithLogContext(logContext)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return withLogContext(svc.WithLog(next))
	}
}

// Welcome method welcomes users
func (s *Service) Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome Odor Profile!!!")
}

// UpdateUserProfile updates profile
func (s *Service) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {

	var request odor.Profile
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		invalidRequestError := svc.NewInvalidRequestError("Bad Request", "Bad Request")
		svc.ReplyWithError(w, r, invalidRequestError)
		return
	}
	s.users <- &request
}

// GetUserProfile returns the user profile
func (s *Service) GetUserProfile(msisdn string) (*odor.Profile, error) {
	if v, ok := profiles[msisdn]; ok {
		p := v
		return p, nil
	}
	return nil, svc.NotFoundError
}

func (s *Service) run() {
	// This function manages all the events that happen in a room
	for {
		select {
		case p := <-s.users:
			profiles[p.MSISDN] = p
			// default:
			// 	fmt.Print("hola")
			// 	// 	close(s.users)
		}
	}
}
