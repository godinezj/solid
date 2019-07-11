# Solid

The solidly.io API.

## Database Setup

...

## Docker

To create ubuntu docker image with above
docker build -t godinezj/solid-api .
docker tag godinezj/solid-api godinezj/solid-api:0.1
docker push godinezj/solid-api

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

`buffalo dev`

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

## LDAP

Set up LDAP container, see [LDAP.md](LDAP.md)

## VPN

Setup the VPN container, see [VPN.md](VPN.md)

## API

Create user:

```bash
curl -H 'Content-Type: application/json' \
    -d '{"username": "JOHNDOE", "first_name": "John", "last_name": "Doe", "email":"jd@example.com", "password":"P@ssw0rd!", "password_confirm":"P@ssw0rd!", "zip": "90210"}' \
    -X POST http://127.0.0.1:3000/users
```

Login

```bash
curl -H 'Content-Type: application/json' \
    -d '{"email":"johndoe", "password":"password"}' \
    -v http://127.0.0.1:3000/login
```

Forgot password

```bash
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com"}' \
    http://127.0.0.1:3000/forgot_password
```

Reset password

```bash
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com", "reset_token_confirm": "507f6c6c-19ca-48a2-9ca8-30f4901e8345", "password":"password1", "password_confirm":"password1"}' \
    http://127.0.0.1:3000/reset_password
```

Login with new password

```bash
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com", "password":"password1"}' \
    http://127.0.0.1:3000/login
```