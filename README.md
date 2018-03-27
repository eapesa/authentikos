Authentikos
-----------

A simple web service for showcasing `generation` (via QR codes) and `verification`
of passcodes produced in TOTP algorithm. These passcodes are used for 2FAs
(Two-Factor Authentication) that provides additional layer of security. Popular
client for TOTP generators is Google Authenticator.

Backend API written in [Go](https://golang.org/).
Frontend managed by [Vue.js](https://vuejs.org/v2/guide/#Getting-Started).


# Dependencies

* [OTP](https://github.com/pquerna/otp)
* [Go-Redis](https://github.com/go-redis/redis)
* Vue.js resource files (already included in `priv/assets/js/ext`)


# Installation

- Ensure `go` is properly configured especially `$GOROOT` and `$GOPATH`.
- Install `go` dependencies:

```
$> go get github.com/pquerna/otp
$> go get github.com/go-redis/redis
```

- Ensure `redis` is started and configure redis credentials
- Install via `make`

```
$> make install
```

*NOTE: You may cleanup your builds by executing `make clean`*

- Default endpoint is `http://localhost:8000`
