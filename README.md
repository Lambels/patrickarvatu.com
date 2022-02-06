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
- FRONTEND
- Near future: Support more OAuth providers
- Write tests: [sql package](https://github.com/Lambels/patrickarvatu.com/tree/master/sqlite)
- Write tests: [http package](https://github.com/Lambels/patrickarvatu.com/tree/master/http)

# Backend
- Powered by Golang
- Architecture layout from [Ben Johnsons Article](https://github.com/benbjohnson/wtf)
- Event Service powered by [asynq](https://github.com/hibiken/asynq)
- CLI start upp powered by [cobra](https://github.com/spf13/cobra)

# Frontend
- Not implemented