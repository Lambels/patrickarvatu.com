# patrickarvatu.com [![Version](https://img.shields.io/badge/goversion-1.17.x-blue.svg)](https://golang.org) ![Go Backend Test](https://github.com/Lambels/patrickarvatu.com/workflows/Go%20Test%20&%20Build/badge.svg)
patrickarvatu.com is my open source personal website!

# State
patrickarvatu.com isnt yet in production, however it is under active development.

### Things Done:
- SQL logic implemented.
- Implement sql code in [sql package](https://github.com/Lambels/patrickarvatu.com/tree/master/sqlite)
- HTTP exposure to the [sql package](https://github.com/Lambels/patrickarvatu.com/tree/master/sqlite)
- OAuth github implementation
- Event Service implemented using [asynq](https://github.com/hibiken/asynq)
- CLI start upp

### TODO:
- Dockerize app
- Node api to interact with frontend fs
- Finish static pages on frontend (about, index)
- Profile component
- 
- Near future: Support more OAuth providers
- Write tests: [sql package](https://github.com/Lambels/patrickarvatu.com/tree/master/sqlite)
- Write tests: [http package](https://github.com/Lambels/patrickarvatu.com/tree/master/http)

# Backend
- Powered by Golang
- Architecture layout from [Ben Johnsons Article](https://github.com/benbjohnson/wtf)
- Event Service powered by [asynq](https://github.com/hibiken/asynq)
- CLI start upp powered by [cobra](https://github.com/spf13/cobra)

# Frontend
- Built with react.js
- Using next.js on top of react.js
- Styling done with tailwind css

# Config File
The config file uses the [toml](https://github.com/toml-lang/toml) formant.
## Fields:
| Field      | Description | Under          |
| :---        |    :----:   |          ---: |
| client-id | Client ID of github oath 2.0 app | [github] |
| client-secret | Client Secret of github oauth 2.0 app | [github] |
| addr | the address of the server (specify only port in development) | [http]
| domain | the domain of the server (leave this empty in development) | [http]
| block-key | key used for secure cookie encryption ([see more](https://github.com/gorilla/securecookie#examples)) | [http]
| hash-key | key used for secure cookie encryption ([see more](https://github.com/gorilla/securecookie#examples)) | [http]
| frontend-url | URL to frontend (ex: http://localhost:3000) | [http]
| admin-user-email | the email of the admin, used to recognize admin user | [user]
| sqlite-dsn | path to sqlite database | [database]
| redis-dsn | redis data source name (ex: 127.0.0.1:6379) | [database]
| addr | address of the smtp server | [smtp]
| identity | refer: [godoc](https://pkg.go.dev/net/smtp#PlainAuth) | [smtp]
| username | refer: [godoc](https://pkg.go.dev/net/smtp#PlainAuth) | [smtp]
| password | refer: [godoc](https://pkg.go.dev/net/smtp#PlainAuth) | [smtp]
| host | refer: [godoc](https://pkg.go.dev/net/smtp#PlainAuth) | [smtp]

# Run:
Currently the app isnt dockerized but you can run the go backend using go command line tool.
```
go install github.com/Lambels/patrickarvatu.com/cmd
```
If you have your GOBIN set to your path run the installed binary with the serve sub command and --config flag
```
bin_name serve --config ./path/to/config/file.toml
```
After you should have a running server on the address and domain specified in the config file.