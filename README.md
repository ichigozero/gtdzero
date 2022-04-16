# gtdzero

ToDo list API built in Go with Gin web framework.

## System Diagram

<p align="center">
  <img src="./assets/diagram.jpg" />
</p>

## API List

### Auth

| HTTP Method | URI      |
| ----------- | -------- |
| POST        | /login   |
| POST        | /logout  |
| POST        | /refresh |

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

### Refresh Token

```bash
curl -i -X "POST" "http://localhost:8080/refresh" \
	-H 'Accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyOTJhOWUwLTAzNjItNDA5My1iZGRkLTk5N2YwM2I5MzljZCIsImV4cCI6MTY1MDcxNDAxOSwicmVmcmVzaF91dWlkIjoiNjdiOTM1MGEtNWQ3Yi01ZjJiLThhMDItZmM0MjgxMzkxZmZjIiwidXNlcl9pZCI6MX0.LW4xBvlFIvL6AIASSyiy6zEW8Mqu-wOf95M7gDcYo2o"}'
```

Output:

```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Sat, 16 Apr 2022 11:42:19 GMT
Content-Length: 496

{"tokens":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTAxMTAyMzksInVzZXJfaWQiOjEsInV1aWQiOiIzYTY3NmQyNi0xY2MyLTQ4NjYtOGEzNC1lYjMzZWI1NGM0MjAifQ.3oRD-vXVrEnA4LG1NjhxJh_-70_I0CGoVZZnAwqRb1w","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cC
I6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjNhNjc2ZDI2LTFjYzItNDg2Ni04YTM0LWViMzNlYjU0YzQyMCIsImV4cCI6MTY1MDcxNDEzOSwicmVmcmVzaF91dWlkIjoiZTc2NjUyMmItMWVmYy01NzY3LWJlODktMWE1NjE0OTJjNjJlIiwidXNlcl9pZCI6MX0.96-3JesxyL_dCZXEH8oml9vGguc3coKomsygUT4RjbA"}}⏎
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
