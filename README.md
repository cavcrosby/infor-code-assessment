# code-assessment

1. How to run the application.

```shell
go mod download
go run ./main.go
```

2. An example cURL command to create a User record via the app’s JSON API.

```shell
curl --request POST localhost:8080/users --header 'Content-Type: application/json' --data '{"id":11, "email":"user@gmail.com", "first_name": "user", "last_name": "user"}'
```

```shell
# additional curl examples
# curl examples to fetch users, valid users, and a invalid user
curl localhost:8080/users
curl localhost:8080/users/1
curl localhost:8080/users/2

# invalid user (well, at least with only the initial populated database and
# newly created user "11")
curl localhost:8080/users/12

# curl examples to update a user, and delete a user
curl --request POST localhost:8080/users/11 --header 'Content-Type: application/json' --data '{"email":"updated_user@gmail.com", "first_name": "user", "last_name": "user"}'
curl --request DELETE localhost:8080/users/11

# curl examples to demonstrate pagination
curl 'localhost:8080/users?page=0&per=5'
curl 'localhost:8080/users?page=1&per=5'

# curl example for pagination with sort/order
curl 'localhost:8080/users?page=0&per=5&sort=id&order=desc'
```

3. Approximately how much time you spent on the application.

**I have spent approximately 4 hours and 31 minutes on this application.**

4. Any tradeoffs you made during development of the application.

* I chose to use Sqlite as that is what I am familiar with using. However, for large read and write operations, I don't believe Sqlite is known to be best tool to use for this kind of requirement. Perfectly fine for a demo though.
* I chose to use the Gin framework for developing the API because it is a framework I have a bit of familiarity with. Beego and Echo may have worked for this API as well.
* A lot of the queries made to the Sqlite database rely on trusting the user to pass in "valid" data (e.g. user id, json data) to the API. This isn't ideal because malicious input could be passed in todo nefarious things to the API. Going past a demo, I would work sanitizing inputs to the API and use prepared queries for reusability.
* For the purposes of the demo, I chose not to create a makefile to include along with the API. To make the project easier to maintain over time (e.g for building, testing, and packaging purposes), a makefile (or something other build tool) would need to be integrated in with the project.
* Timestamps are handled solely on the database side. The API does not generate a timestamp, attach it to a newly created (or updated) user, and pass this to the database. Rather, I decide to let the database handle this by itself and the API only sees this later when retrieving a user. I did this in the interest of speed.
