package main

import (
	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/asynq"
	"github.com/Lambels/patrickarvatu.com/http"
	"github.com/Lambels/patrickarvatu.com/smtp"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

func newDB(cfg *pa.Config) (*sqlite.DB, func(), error) {
	db := sqlite.NewDB(cfg.Database.SqliteDSN)

	if err := db.Open(); err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
	}, nil
}

func newEventService(cfg *pa.Config) (*asynq.EventService, func(), error) {
	eService := asynq.NewEventService(cfg.Database.RedisDSN)

	return eService, func() {
		eService.Close()
	}, nil
}

func newEmailService(cfg *pa.Config) pa.EmailService {
	return smtp.NewEmailService(cfg.Smtp.Addr, cfg.Smtp.Identity, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Host)
}

func newServer(cfg *pa.Config,
	authService pa.AuthService,
	userService pa.UserService,
	blogService pa.BlogService,
	subBlogService pa.SubBlogService,
	commentService pa.CommentService,
	eventService *asynq.EventService,
	subscriptionService pa.SubscriptionService,
	emailService pa.EmailService,
) (*http.Server, func(), error) {
	s := http.NewServer(cfg)

	s.AuthService = authService
	s.UserService = userService
	s.BlogService = blogService
	s.SubBlogService = subBlogService
	s.CommentService = commentService
	s.EventService = eventService
	s.SubscriptionService = subscriptionService
	s.EmailService = emailService

	s.EventService.RegisterSubscriptionsHandler(s.SubscriptionService)
	s.EventService.RegisterHandler(pa.EventTopicNewComment, s.HandleCommentEvent)
	s.EventService.RegisterHandler(pa.EventTopicNewSubBlog, s.HandleSubBlogEvent)

	// open registered event service.
	if err := eventService.Open(); err != nil {
		return nil, nil, err
	}

	s.Addr = cfg.HTTP.Addr
	s.Domain = cfg.HTTP.Domain

	if err := s.Open(); err != nil {
		return nil, nil, err
	}

	return s, func() {
		s.Close()
	}, nil
}

func initializeServer(cfg *pa.Config) (*http.Server, func(), error) {
	db, clnUpDB, err := newDB(cfg)
	if err != nil {
		return nil, nil, err
	}

	evSrv, clnUpEvSrv, err := newEventService(cfg)
	if err != nil {
		clnUpDB()
		return nil, nil, err
	}

	emSrv := newEmailService(cfg)

	auSrv := sqlite.NewAuthService(db)
	usSrv := sqlite.NewUserService(db)
	blSrv := sqlite.NewBlogService(db)
	sbSrv := sqlite.NewSubBlogService(db)
	cmSrv := sqlite.NewCommentService(db)
	subSrv := sqlite.NewSubscriptionService(db)

	serv, clnUpServ, err := newServer(
		cfg,
		auSrv,
		usSrv,
		blSrv,
		sbSrv,
		cmSrv,
		evSrv,
		subSrv,
		emSrv,
	)
	if err != nil {
		clnUpDB()
		clnUpEvSrv()
	}

	return serv, func() {
		clnUpDB()
		clnUpEvSrv()
		clnUpServ()
	}, nil
}
