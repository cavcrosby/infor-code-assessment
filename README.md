# code-assessment

1. How to run the application.

```shell
go mod download
go run ./main.go
```

2. An example cURL command to create a User record via the app’s JSON API.

```shell
# curl examples to fetch users, valid users, and a invalid user
curl localhost:8080/users
curl localhost:8080/users/1
curl localhost:8080/users/2

# invalid user (well, at least with only the initial populated database)
curl localhost:8080/users/11

# curl examples to create a user, update a user, and delete a user
curl --request POST localhost:8080/users --header 'Content-Type: application/json' --data '{"id":11, "email":"f@gmail.com", "first_name": "foo", "last_name": "foo"}'
curl localhost:8080/users/11

curl --request POST localhost:8080/users/11 --header 'Content-Type: application/json' --data '{"email":"d@gmail.com", "first_name": "foo", "last_name": "foo"}'
curl localhost:8080/users/11
curl --request DELETE localhost:8080/users/11

# curl example to demonstrate pagination
curl 'localhost:8080/users?page=0&per=5'
curl 'localhost:8080/users?page=1&per=5'
```

3. Approximately how much time you spent on the application.

**I have spent approximately 4 hours on this application.**

4. Any tradeoffs you made during development of the application.

* I chose to use Sqlite as that is what I am familiar with using. However, for large read and write operations, I don't believe Sqlite is known to be best tool to use for this kind of requirement. Perfectly fine for a demo though.
* I chose to use the Gin framework for developing the API because Golang shows off a tutorial for using this framework to make an API. Having followed this tutorial before, I decided to utilize it in getting a barebones idea of how to write this REST API.
* A lot of the queries made to the Sqlite database rely on trusting the user to pass in "valid" data (e.g. user id, json data) to the API. This is a horrible idea because malicious input could be passed in todo nefarious things to the API. Going past a demo, I would work sanitizing inputs to the API and use prepared queries for reusability.
* For the purposes of the demo, I chose not to create a makefile to include along with the API. To make the project easier to maintain over time (e.g for building, testing, and packaging purposes), a makefile (or something other build tool) would need to be integrated in with the project.
