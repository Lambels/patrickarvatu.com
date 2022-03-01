package main

import (
	"log"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/asynq"
	"github.com/Lambels/patrickarvatu.com/fs"
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

func newFileService(root string) pa.FileService {
	return fs.NewFileService(root)
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
	projectService pa.ProjectService,
	projectsFileSystem pa.FileService,
	blogsFileSystem pa.FileService,
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
	s.ProjectService = projectService
	s.ProjectsFileSystem = projectsFileSystem
	s.BlogsFileSystem = blogsFileSystem

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
	log.Println("[DEBUG] Connected to db.")

	evSrv, clnUpEvSrv, err := newEventService(cfg)
	if err != nil {
		clnUpDB()
		return nil, nil, err
	}
	log.Println("[DEBUG] Initialized event service.")

	emSrv := newEmailService(cfg)
	log.Println("[DEBUG] Initialized email service.")

	prFs := newFileService("/images/projects")
	log.Println("[DEBUG] Initialized projects images file system.")

	blFs := newFileService("/images/blogs")
	log.Println("[DEBUG] Initialized blogs images file system.")

	auSrv := sqlite.NewAuthService(db)
	usSrv := sqlite.NewUserService(db)
	blSrv := sqlite.NewBlogService(db)
	sbSrv := sqlite.NewSubBlogService(db)
	cmSrv := sqlite.NewCommentService(db)
	subSrv := sqlite.NewSubscriptionService(db)
	pjSrv := sqlite.NewProjectService(db)
	log.Println("[DEBUG] Started database services.")

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
		pjSrv,
		prFs,
		blFs,
	)
	if err != nil {
		clnUpDB()
		clnUpEvSrv()
	}
	log.Println("[INFO] Started server on address", serv.Addr)

	return serv, func() {
		clnUpDB()
		clnUpEvSrv()
		clnUpServ()
	}, nil
}
