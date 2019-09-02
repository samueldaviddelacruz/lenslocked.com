# Lens Locked
A photo gallery application written in Go.

## How to Run the project
Install dependencies (go 1.11+ required):
```
go get
```
Then proceed to build and run the project
```
go run *.go
```
it will be running on localhost:4000

## Built With

* [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) - For http routing
* [Gorilla CSRF](gorilla/csrf) - For Cross Site Request Forgery (CSRF) prevention
* [GORM](https://gorm.io/) - For database interaction
* [Postgress](https://www.postgresql.org/) - For persistency
* [Mailgun](https://www.mailgun.com/) - For sending emails
* [UNOFFICIAL Dropbox Go SDK](https://github.com/dropbox/dropbox-sdk-go-unofficial) - For dropbox interaction
* [oauth2 package](https://godoc.org/golang.org/x/oauth2) - For authenticating with Oauth2 services
* [Digital Ocean](https://www.digitalocean.com) - For deployment.
* [Caddy Server](https://caddyserver.com/) - For HTTP proxy and sane security defaults.

## Demo
* You can try it out here => [Demo](https://lenslocked-project-demo.net/) 