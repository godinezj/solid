# Welcome to Buffalo!

Thank you for choosing Buffalo for your web development needs.

## Database Setup

It looks like you chose to set up your application using a postgres database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start postgres for you.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started postgres, now Buffalo can create the databases in that file for you:

	$ buffalo db create -a

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

	$ buffalo dev

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

## API

Create user:

```
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com", "password":"password", "password_confirm":"password"}' \
    -X POST http://127.0.0.1:3000/users
```

Login

```
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com", "password":"password"}' \
    http://127.0.0.1:3000/login
```

Forgot password

```
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com"}' \
    http://127.0.0.1:3000/forgot_password
```

Reset password

```
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com", "reset_token_confirm": "507f6c6c-19ca-48a2-9ca8-30f4901e8345", "password":"password1", "password_confirm":"password1"}' \
    http://127.0.0.1:3000/reset_password
```

Login with new password

```
curl -H 'Content-Type: application/json' \
    -d '{"email":"godinezj@gmail.com", "password":"password1"}' \
    http://127.0.0.1:3000/login
```