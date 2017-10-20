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
var radiusMap map[string]*odor.RadiusPacket

// Service represents the Profile service
type Service struct {
	*odor.Config
	server     *http.Server
	usersChan  chan *odor.Profile
	radiusChan chan *odor.RadiusPacket
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
	s.usersChan = make(chan *odor.Profile)
	profiles = map[string]*odor.Profile{}
	radiusMap = map[string]*odor.RadiusPacket{}
	go s.run()
	return s.server.ListenAndServe()
}

// Stop the service
func (s *Service) Stop() error {
	if s.server != nil {
		close(s.usersChan)
		s.server.Shutdown(nil)
		s.server = nil
	}
	return nil
}

func (s *Service) router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/{msisdn}", s.withMws("updateUserProfile")(s.UpdateUserProfile)).Methods("PUT")
	r.HandleFunc("/ips/{ip}", s.withMws("InjectRadiusPacket")(s.InjectRadiusPacket)).Methods("PUT")
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

// InjectRadiusPacket injects a radius packet to map (ip, msisdn)
func (s *Service) InjectRadiusPacket(w http.ResponseWriter, r *http.Request) {
	var request odor.RadiusPacket
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		invalidRequestError := svc.NewInvalidRequestError("Bad Request", "Bad Request")
		svc.ReplyWithError(w, r, invalidRequestError)
		return
	}

	// netRadiusP, err := odor.RadiusPacket2NetRadiusPacket(request)
	// if err != nil {
	// 	svc.ReplyWithError(w, r, err)
	// 	return
	// }
	s.radiusChan <- &request
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
	s.usersChan <- &request
}

// GetUserProfile returns the user profile
func GetUserProfile(msisdn string) (*odor.Profile, error) {
	svc.NewLogger().Info("GetUserProfile for msisdn: %s", msisdn)
	if v, ok := profiles[msisdn]; ok {
		svc.NewLogger().Info("Obtained profile: %+v", v)
		return v, nil
	}
	svc.NewLogger().Warn("No profile for msisdn %s", msisdn)
	return nil, svc.NotFoundError
}

// GetRadiusPacket returns the radius packet associated to an IP
func GetRadiusPacket(ip string) (*odor.RadiusPacket, error) {
	svc.NewLogger().Info("GetRadiusPacket for IP: %s", ip)
	if v, ok := radiusMap[ip]; ok {
		svc.NewLogger().Info("Obtained mapping: %+v", v)
		return v, nil
	}
	svc.NewLogger().Warn("No mapping for IP %s", ip)
	return nil, svc.NotFoundError
}

func (s *Service) run() {
	// This function manages all the events that happen in a room
	for {
		select {
		case p := <-s.usersChan:
			profiles[p.MSISDN] = p
		case r := <-s.radiusChan:
			radiusMap[r.IP] = r
			// default:
			// 	fmt.Print("hola")
			// 	// 	close(s.users)
		}
	}
}
