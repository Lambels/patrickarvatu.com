package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/securecookie"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// TODO: publish event after sql creation

// ServerShutdownTime is the time the server allows processes to finish before shuting down
const ServerShutdownTime = 3 * time.Second

// Server represents a backend api HTTP service which wraps all our HTTP functionality.
type Server struct {
	server *http.Server
	router *chi.Mux
	ln     net.Listener
	sc     *securecookie.SecureCookie

	// server address.
	Addr   string
	Domain string

	// Services used by the http package.
	AuthService      pa.AuthService
	UserService      pa.UserService
	BlogService      pa.BlogService
	SuBlogService    pa.SubBlogService
	CommentService   pa.CommentService
	EventService     pa.EventService
	SubsctionService pa.SubscriptionService

	conf *pa.Config
}

// NewServer registers all routes and returns a new server with the conf configurations.
func NewServer(conf *pa.Config) *Server {
	s := &Server{
		server: &http.Server{},
		router: chi.NewRouter(),
		conf:   conf,
	}

	// middleware stack.
	s.router.Use(chimw.Logger)
	s.router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{s.conf.HTTP.FrontendURL},
			AllowedMethods:   []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodOptions, http.MethodPut},
			AllowCredentials: true,
		},
	))

	// set custom not found api handler.
	s.router.NotFound(s.handleNotFound)

	s.router.Route("/oauth", func(r chi.Router) {
		s.registerAuthRoutes(r)
	})

	// register router to server with registered routes.
	s.server.Handler = s.router

	return s
}

// Open validates the configurration and starts the server on the address.
func (s *Server) Open() error {
	// validate our providers keys.
	if s.conf.Github.ClientID == "" {
		return fmt.Errorf("gtihub client id not set.")
	} else if s.conf.Github.ClientSecret == "" {
		return fmt.Errorf("github client secret not set.")
	}

	// open the secure cookie implementation.
	if err := s.openSecureCookie(); err != nil {
		return err
	}

	// set listener.
	if s.Domain != "" {
		s.ln = autocert.NewListener(s.Domain)
	} else {
		ln, err := net.Listen("tcp", s.Addr)
		if err != nil {
			return err
		}

		s.ln = ln
	}

	// start listening with s.ln
	go s.server.Serve(s.ln)

	return nil
}

// Close brings the server to a gracefull shutdown.
func (s *Server) Close() error {
	cancelCtx, cancel := context.WithTimeout(context.Background(), ServerShutdownTime)
	defer cancel() // release resources.
	return s.server.Shutdown(cancelCtx)
}

// openSecureCookie uses the keys under the config and checks their existance.
func (s *Server) openSecureCookie() error {
	// ensure keys are set.
	if s.conf.HTTP.BlockKey == "" {
		return fmt.Errorf("block key isnt set.")
	} else if s.conf.HTTP.HashKey == "" {
		return fmt.Errorf("hash key isnt set.")
	}

	s.sc = securecookie.New([]byte(s.conf.HTTP.HashKey), []byte(s.conf.HTTP.BlockKey))
	s.sc.SetSerializer(securecookie.JSONEncoder{}) // use the json encoder.
	return nil
}

// RunDebugServer runs a debug server on port 8000
// blocking function.
func RunDebugServer() error {
	s := http.NewServeMux()
	s.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":8000", s)
}

// NewOAuthConfig returns an oauth2.0 config object to start the oauth2.0 authorization flow.
// source determines the provider config.
// returns an empty config if not found but that should never happen.
func (s *Server) NewOAuthConfig(source string) *oauth2.Config {
	switch source {
	case pa.AuthSourceGitHub:
		return &oauth2.Config{
			ClientID:     s.conf.Github.ClientID,
			ClientSecret: s.conf.Github.ClientSecret,
			Endpoint:     github.Endpoint,
		}

	default:
		return &oauth2.Config{}
	}
}

// handleNotFound sends a not found error with the path.
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	SendError(w, r, pa.Errorf(pa.ENOTFOUND, "%s didnt match with any path.", r.URL.Path))
}

// cookie geter and seter ----------------------------------------------------

// getSession returns a pa.Session from r.
// returns an empty session on error.
func (s *Server) getSession(r *http.Request) (ses *pa.Session, err error) {
	c, err := r.Cookie(pa.SessionCookieName)
	if err != nil {
		return ses, err
	}

	if err := s.sc.Decode(pa.SessionCookieName, c.Value, ses); err != nil {
		return ses, err
	}
	return ses, nil
}

// setSession sets ses on w using the securecookie service.
func (s *Server) setSession(w http.ResponseWriter, ses *pa.Session) error {
	v, err := s.sc.Encode(pa.SessionCookieName, ses)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     pa.SessionCookieName,
		Value:    v,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	return nil
}

// Middleware -----------------------------------------------------------------------

// jsonResponseTypeMiddleware sets the Content-Type header to application/json.
func (s *Server) jsonResponseTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// authentificateMiddleware authentificates a requests based on api key or cookie.
func (s *Server) authentificateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check fo api key.
		if header := r.Header.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
			apiKey := strings.TrimPrefix(header, "Bearer ")

			// api key filter.
			filter := pa.UserFilter{
				APIKey: &apiKey,
			}

			// look for user paired with api key.
			users, _, err := s.UserService.FindUsers(r.Context(), filter)
			if err != nil {
				SendError(w, r, err)
				return
			} else if len(users) == 0 {
				SendError(w, r, pa.Errorf(pa.EUNAUTHORIZED, "api key invalid"))
			}

			// set auth user to ctx and dispatch next handler.
			r = r.WithContext(pa.NewContextWithUser(r.Context(), users[0]))
			next.ServeHTTP(w, r)
			return
		}

		ses, _ := s.getSession(r)

		// if we have a user id under session we fetch user.
		if ses.UserID != 0 {
			if user, err := s.UserService.FindUserByID(r.Context(), ses.UserID); err != nil {
				log.Printf("FindUserByID: id=%v err=%s", ses.UserID, err)
			} else { // user found, ok.
				r = r.WithContext(pa.NewContextWithUser(r.Context(), user))
			}
		}

		next.ServeHTTP(w, r)
	})
}

// adminAuthMiddleware rejects all requests from not admin users.
func (s *Server) adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send EUNAUTHORIZED if ctx isnt admin ctx.
		if pa.IsAdminContext(r.Context()) {
			next.ServeHTTP(w, r)
			return
		}

		SendError(w, r, pa.Errorf(pa.EUNAUTHORIZED, "user is not admin."))
	})
}

// requireAuthMiddleware rejects all requests from non auth users.
func (s *Server) requireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send EUNAUTHORIZED if ctx isnt auth ctx.
		if userID := pa.UserIDFromContext(r.Context()); userID != 0 {
			next.ServeHTTP(w, r)
			return
		}

		SendError(w, r, pa.Errorf(pa.EUNAUTHORIZED, "user is not auth."))
	})
}

// requireNoAuthMiddleware rejects all requests from auth users.
func (s *Server) requireNoAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send EUNAUTHORIZED if ctx is auth ctx.
		if userID := pa.UserIDFromContext(r.Context()); userID == 0 {
			next.ServeHTTP(w, r)
			return
		}

		SendError(w, r, pa.Errorf(pa.EUNAUTHORIZED, "user is auth."))
	})
}
