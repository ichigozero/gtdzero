# gtdzero

ToDo list API built in Go with Gin web framework.

## API List

### Auth

| HTTP Method | URI     |
| ----------- | ------- |
| POST        | /login  |
| POST        | /logout |

### Task

| HTTP Method | URI                     |
| ----------- | ----------------------- |
| GET         | /todo/api/v1.0/tasks    |
| GET         | /todo/api/v1.0/task/:id |
| POST        | /todo/api/v1.0/tasks    |
| PUT         | /todo/api/v1.0/task/:id |
| DELETE      | /todo/api/v1.0/task/:id |

## Examples

### Login

```bash
$ curl -i -X "POST" "http://localhost:8080/login" \
	-H 'Accept: application/json' \
	-H 'Content-Type: application/json' \
	-d $'{"username": "admin", "password": "password"}'
```

Output:

```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Thu, 14 Apr 2022 12:26:51 GMT
Content-Length: 439

{"tokens":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NDk5NDAxMTEsInVzZXJfaWQiOjEsInV1aWQiOiJiY2IwZjY4ZS03ZjQwLTQ3ZmYtYTM3OS0wZmRlOGEzNDBkNzUifQ.biW4Iz1JKrPdHdmLxoR5Z2VXsbFY9GyvuGXCDJuvRqo","refresh_token":"
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTA1NDQwMTEsInVzZXJfaWQiOjEsInV1aWQiOiIyY2UzZTlkYS1kZGE5LTQ3OTMtYTE1Ni04ZjMwMzM5ZWQxM2UifQ.pP6RPF3g8NJvmc6ihm-zWH2if3Oz7XL7Ci957ekicbM"}}⏎
```

### Logout

```bash
curl -i -X "POST" "http://localhost:8080/logout" \
	-H 'Accept: application/json' \
	-H 'Content-Type: application/json' \
	-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NDk5NDAxMTEsInVzZXJfaWQiOjEsInV1aWQiOiJiY2IwZjY4ZS03ZjQwLTQ3ZmYtYTM3OS0wZmRlOGEzNDBkNzUifQ.biW4Iz1JKrPdHdmLxoR5Z2VXsbFY9GyvuGXCDJuvRqo'
```

Output:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 14 Apr 2022 12:30:31 GMT
Content-Length: 15

{"result":true}
```

### Create Task

```bash
$ curl -i -X "POST" "http://localhost:8080/todo/api/v1.0/tasks" \
	-H 'Accept: application/json' \
	-H 'Content-Type: application/json' \
	-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NDk5NDAxMTEsInVzZXJfaWQiOjEsInV1aWQiOiJiY2IwZjY4ZS03ZjQwLTQ3ZmYtYTM3OS0wZmRlOGEzNDBkNzUifQ.biW4Iz1JKrPdHdmLxoR5Z2VXsbFY9GyvuGXCDJuvRqo' \
	-d '{"title":"Read a book", "description": "Foo"}'
```

Output:

```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Thu, 14 Apr 2022 12:35:31 GMT
Content-Length: 83

{"task":{"id":3,"title":"Read a book","description":"Foo","done":false,"userid":1}}
```
