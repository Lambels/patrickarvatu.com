package http

import (
	"context"
	"encoding/json"
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
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// ServerShutdownTime is the time the server allows processes to finish before shuting down
const ServerShutdownTime = 3 * time.Second

// Server represents a backend api HTTP service which wraps all our HTTP functionality.
type Server struct {
	server *http.Server
	router *chi.Mux
	ln     net.Listener
	sc     *securecookie.SecureCookie
	cron   *cron.Cron

	// server address.
	Addr   string
	Domain string

	// Services used by the http package.
	AuthService         pa.AuthService
	UserService         pa.UserService
	BlogService         pa.BlogService
	SubBlogService      pa.SubBlogService
	CommentService      pa.CommentService
	EventService        pa.EventService
	SubscriptionService pa.SubscriptionService
	EmailService        pa.EmailService

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
	s.router.Use(s.authentificateMiddleware) // attempt to authentificate each request.
	s.router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{s.conf.HTTP.FrontendURL},
			AllowedMethods:   []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodOptions, http.MethodPut},
			AllowCredentials: true,
		},
	))

	// set custom not found api handler.
	s.router.NotFound(s.handleNotFound)

	s.router.Route("/v1/oauth", func(r chi.Router) {
		s.registerAuthRoutes(r)
	})

	s.router.Route("/v1/users", func(r chi.Router) {
		r.Use(s.requireAuthMiddleware)
		s.registerUserRoutes(r)
	})

	s.router.Route("/v1/blogs", func(r chi.Router) {
		s.registerBlogRoutes(r)
	})

	s.router.Route("/v1/sub-blogs", func(r chi.Router) {
		s.registerSubBlogRoutes(r)
	})

	s.router.Route("/v1/comments", func(r chi.Router) {
		r.Use(s.requireAuthMiddleware)
		s.registerCommentRoutes(r)
	})

	s.router.Route("/v1/subscriptions", func(r chi.Router) {
		r.Use(s.requireAuthMiddleware)
		s.registerSubscriptionRoutes(r)
	})

	// register router to server with registered routes.
	s.server.Handler = s.router

	return s
}

// Open validates the configurration and starts the server on the address.
func (s *Server) Open() error {
	// validate our providers keys.
	if s.conf.Github.ClientID == "" {
		return fmt.Errorf("gtihub client id not set")
	} else if s.conf.Github.ClientSecret == "" {
		return fmt.Errorf("github client secret not set")
	}

	// open the secure cookie implementation.
	if err := s.openSecureCookie(); err != nil {
		return err
	}

	// open cronjob.
	s.openCronJob()

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

// UseTLS retruns true if cert + key file found.
func (s *Server) UseTLS() bool {
	return s.Domain != ""
}

// Close brings the server to a gracefull shutdown and stops the cron job.
func (s *Server) Close() error {
	s.cron.Stop() // stop the cron job.

	cancelCtx, cancel := context.WithTimeout(context.Background(), ServerShutdownTime)
	defer cancel() // release resources.
	return s.server.Shutdown(cancelCtx)
}

// openSecureCookie uses the keys under the config and checks their existance.
func (s *Server) openSecureCookie() error {
	// ensure keys are set.
	if s.conf.HTTP.BlockKey == "" {
		return fmt.Errorf("block key isnt set")
	} else if s.conf.HTTP.HashKey == "" {
		return fmt.Errorf("hash key isnt set")
	}

	s.sc = securecookie.New([]byte(s.conf.HTTP.HashKey), []byte(s.conf.HTTP.BlockKey))
	s.sc.SetSerializer(securecookie.JSONEncoder{}) // use the json encoder.
	return nil
}

// cronJob ------------------------------------------------------------------

// openCronJob creates a new cron job with cron.DefaultLogger.
func (s *Server) openCronJob() {
	s.cron = cron.New(cron.WithLogger(
		cron.DefaultLogger,
	))
	s.cron.Start()
}

// RegisterCronJob registers a new job to s.cron, can be call regardless if cron is or isnt running.
func (s *Server) RegisterCronJon(spec string, cmd func()) error {
	_, err := s.cron.AddFunc(spec, cmd)
	return err
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

// publishNewEvent is a helper function to push an event on the event queue.
func (s *Server) publishNewEvent(ctx context.Context, event pa.Event) error {
	return s.EventService.Push(ctx, event)
}

// handleNotFound sends a not found error with the path.
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	SendError(w, r, pa.Errorf(pa.ENOTFOUND, "%s didnt match with any path.", r.URL.Path))
}

// cookie geter and seter ----------------------------------------------------

// getSession returns a pa.Session from r.
// returns an empty session on error.
func (s *Server) getSession(r *http.Request) (pa.Session, error) {
	c, err := r.Cookie(pa.SessionCookieName)
	if err != nil { // simply return an empty session if cookie isnt present.
		return pa.Session{}, nil
	}

	var ses pa.Session
	if err := s.sc.Decode(pa.SessionCookieName, c.Value, &ses); err != nil {
		return pa.Session{}, err
	}
	return ses, nil
}

// setSession sets ses on w using the securecookie service.
func (s *Server) setSession(w http.ResponseWriter, ses pa.Session) error {
	v, err := s.sc.Encode(pa.SessionCookieName, ses)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     pa.SessionCookieName,
		Value:    v,
		Path:     "/",
		Secure:   s.UseTLS(),
		Domain:   "localhost", // pass cookie to all sub domains including frontend and api.
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

// Event Handlers -----------------------------------------------------------------

// HandleCommentEvent handels the pa.EventTopicNewComment -> ./event.go.
// sends an email to all subscribers.
func (s *Server) HandleCommentEvent(ctx context.Context, hand pa.SubscriptionService, event pa.Event) error {
	var payload pa.CommentPayload
	if err := json.Unmarshal(event.Payload.([]byte), &payload); err != nil {
		log.Println("[UnMarshalError] err: ", err.Error())
		return err
	}

	v := pa.EventTopicNewComment
	subs, _, err := hand.FindSubscriptions(ctx, pa.SubscriptionFilter{
		Topic:   &v,
		Payload: payload,
	})
	if err != nil {
		log.Println("[FindSubscriptions] err: ", err.Error())
		return err
	} else if len(subs) == 0 { // no subscriptions.
		return fmt.Errorf("no subscriptions found")
	}

	// add recievers.
	var to []string
	for _, sub := range subs {
		usr, err := s.UserService.FindUserByID(ctx, sub.UserID)
		if err != nil {
			log.Println("[FindUserByID] err: ", err.Error())
			continue
		}

		// if user has attached email add him in to.
		if usr.Email != "" {
			to = append(to, usr.Email)
		}
	}

	// we have reciepients.
	if len(to) != 0 {
		subBlog, err := s.SubBlogService.FindSubBlogByID(ctx, payload.SubBlogID)
		if err != nil {
			log.Println("[FindSubBlogByID] err: ", err.Error())
			return err
		}

		if err := s.EmailService.SendEmail(to,
			fmt.Sprintf("There's been a new comment on %s, go check it out! %s", subBlog.Title, s.conf.HTTP.FrontendURL+"/sub-blog/"+fmt.Sprint(subBlog.ID)),
			fmt.Sprintf("New Comment On %s", subBlog.Title),
		); err != nil {
			log.Println("[SendEmail] err: ", err.Error())
			return err
		}
	}
	return nil
}

// HandleSubBlogtEvent handels the pa.EventTopicNewSubBlog -> ./event.go.
// sends an email to all subscribers.
func (s *Server) HandleSubBlogEvent(ctx context.Context, hand pa.SubscriptionService, event pa.Event) error {
	var payload pa.SubBlogPayload
	if err := json.Unmarshal(event.Payload.([]byte), &payload); err != nil {
		log.Println("[UnMarshalError] err: ", err.Error())
		return err
	}

	v := pa.EventTopicNewSubBlog
	subs, _, err := hand.FindSubscriptions(ctx, pa.SubscriptionFilter{
		Topic:   &v,
		Payload: payload,
	})
	if err != nil {
		log.Println("[FindSubscriptions] err: ", err.Error())
		return err
	} else if len(subs) == 0 { // no subscriptions.
		return fmt.Errorf("no subscriptions found")
	}

	// add recievers.
	var to []string
	for _, sub := range subs {
		usr, err := s.UserService.FindUserByID(ctx, sub.UserID)
		if err != nil {
			log.Println("[FindUserByID] err: ", err.Error())
			continue
		}

		// if user has attached email add him in to.
		if usr.Email != "" {
			to = append(to, usr.Email)
		}
	}

	// we have reciepients.
	if len(to) != 0 {
		blog, err := s.BlogService.FindBlogByID(ctx, payload.BlogID)
		if err != nil {
			log.Println("[FindBlogByID] err: ", err.Error())
			return err
		}

		if err := s.EmailService.SendEmail(to,
			fmt.Sprintf("There's been a new article on %s, go check it out! %s", blog.Title, s.conf.HTTP.FrontendURL+"/blog/"+fmt.Sprint(blog.ID)),
			fmt.Sprintf("New Article On %s", blog.Title),
		); err != nil {
			log.Println("[SendEmail] err: ", err.Error())
			return err
		}
	}
	return nil
}
