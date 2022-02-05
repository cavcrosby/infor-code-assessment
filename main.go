package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiUrl       = "localhost:8080"
	dbFilePath   = "infor.db"
	pageMaxLimit = 10
	tableName    = "table1"
)

type paginationResponse struct {
	ApiUrl  string `json:"base"`
	Next    string `json:"next"`
	Results []user `json:"results"`
}

type user struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Updated   string `json:"updated"`
}

var (
	dbConnector      *sql.DB
	initSqlTableStmt = fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %v(
		id INT NOT NULL,
		email TEXT NULL,
		first_name TEXT NULL,
		last_name TEXT NULL,
		updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
		PRIMARY KEY(id)
	);
	`, tableName)

	initUsers = []user{
		{ID: 1, Email: "foobar@gmail.com", FirstName: "foo", LastName: "bar"},
		{ID: 2, Email: "baz@gmail.com", FirstName: "baz", LastName: "baz"},
		{ID: 3, Email: "alice@gmail.com", FirstName: "alice", LastName: "alice"},
		{ID: 4, Email: "bob@gmail.com", FirstName: "bob", LastName: "bob"},
		{ID: 5, Email: "john@gmail.com", FirstName: "john", LastName: "john"},
		{ID: 6, Email: "conner@gmail.com", FirstName: "conner", LastName: "conner"},
		{ID: 7, Email: "smith@gmail.com", FirstName: "smith", LastName: "smith"},
		{ID: 8, Email: "bazzzz@gmail.com", FirstName: "bazzzz", LastName: "bazzzz"},
		{ID: 9, Email: "outofnames@gmail.com", FirstName: "outofnames", LastName: "outofnames"},
		{ID: 10, Email: "google@gmail.com", FirstName: "google", LastName: "google"},
	}
)

func updateData(sqlStatement string) error {
	_, err := dbConnector.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func retrieveData(sqlStatement string) (*sql.Rows, error) {
	rows, err := dbConnector.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func sqliteUnMarshall(rows *sql.Rows, usersFromDb []user) ([]user, error) {
	for rows.Next() {
		var userInDb user
		err := rows.Scan(&userInDb.ID, &userInDb.Email, &userInDb.FirstName, &userInDb.LastName, &userInDb.Updated)
		if err != nil {
			return nil, err
		}
		usersFromDb = append(usersFromDb, userInDb)
	}

	return usersFromDb, nil
}

// Add a user from JSON received in the request body.
func deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := updateData(
		fmt.Sprintf(
			"DELETE FROM %v WHERE id = %v",
			tableName,
			id,
		),
	); err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

// Add a user from JSON received in the request body.
func updateUser(c *gin.Context) {
	var updateUser user
	id := c.Param("id")

	if err := c.BindJSON(&updateUser); err != nil {
		return
	}
	iid, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	// I do not expect a user to pass in an ID to the json in the request body,
	// however, the response JSON should include the ID of the user updated.
	updateUser.ID = iid

	if err := updateData(
		fmt.Sprintf(
			"UPDATE %v SET email = '%v', first_name = '%v', last_Name = '%v' WHERE id = %v",
			tableName,
			updateUser.Email,
			updateUser.FirstName,
			updateUser.LastName,
			id,
		),
	); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, updateUser)
}

// Add an user from JSON received in the request body.
func createUser(c *gin.Context) {
	var newUser user

	// Call BindJSON to bind the received JSON to
	// newUser.
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	if err := updateData(
		fmt.Sprintf(
			"INSERT INTO %v(id, email, first_name, last_name) values(%v, '%v', '%v', '%v')",
			tableName,
			newUser.ID,
			newUser.Email,
			newUser.FirstName,
			newUser.LastName,
		),
	); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")

	row, err := retrieveData(fmt.Sprintf("SELECT * FROM %v WHERE id = %v", tableName, id))
	if err != nil {
		log.Fatal(err)
	}

	var usersFromDb []user
	usersFromDb, err = sqliteUnMarshall(row, usersFromDb)
	if err != nil {
		log.Fatal(err)
	}

	if usersFromDb != nil {
		c.IndentedJSON(http.StatusOK, usersFromDb[0])
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// Respond with the list of all users as JSON.
func getUsers(c *gin.Context) {
	var paginationPage string
	var paginationPer string
	for k, v := range c.Request.URL.Query() {
		if k == "page" {
			paginationPage = v[0]
		} else if k == "per" {
			paginationPer = v[0]
		}
	}

	var paginationPageI, paginationPerI int
	var offsetPagPageI int = 0
	var err error
	if paginationPage != "" && paginationPer != "" {
		paginationPageI, err = strconv.Atoi(paginationPage)
			if err != nil {
				return
			}

		paginationPerI, err = strconv.Atoi(paginationPer)
				if err != nil {
					return
				}

		// Without this patch, when pages > 0, the last result is included each time.
		// Which is something I do not want.
		if paginationPage != "0" {
			offsetPagPageI += 1
		}
	}

	if paginationPage != "" && paginationPer != "" {
		rows, err := retrieveData(fmt.Sprintf("SELECT * FROM %v WHERE id >= %v LIMIT '%v'", tableName, (paginationPageI*paginationPerI + offsetPagPageI), paginationPerI))
		if err != nil {
			log.Fatal(err)
		}

		var usersFromDb []user
		usersFromDb, err = sqliteUnMarshall(rows, usersFromDb)
		if err != nil {
			log.Fatal(err)
		}

		paginationPerINext, err := strconv.Atoi(paginationPer)
		if err != nil {
			return
		}

		pagiRes := paginationResponse{
			ApiUrl:  apiUrl,
			Next:    fmt.Sprintf(apiUrl+"/users/?page=%v&per=%v", paginationPageI+1, paginationPerINext),
			Results: usersFromDb,
		}

		c.IndentedJSON(http.StatusOK, pagiRes)
	} else {
		rows, err := retrieveData(fmt.Sprintf("SELECT * FROM %v", tableName))
		if err != nil {
			log.Fatal(err)
		}

		var usersFromDb []user
		usersFromDb, err = sqliteUnMarshall(rows, usersFromDb)
		if err != nil {
			log.Fatal(err)
		}

		c.IndentedJSON(http.StatusOK, usersFromDb)
	}
}

func main() {
	var err error
	dbConnector, err = sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnector.Close()

	if _, err := os.Stat(dbFilePath); errors.Is(err, fs.ErrNotExist) {

		_, err = dbConnector.Exec(initSqlTableStmt)
		if err != nil {
			log.Panic(err)
		}

		for _, u := range initUsers {
			_, err = dbConnector.Exec(
				fmt.Sprintf(
					"INSERT INTO %v(id, email, first_name, last_name) values('%v', '%v', '%v', '%v')",
					tableName,
					u.ID,
					u.Email,
					u.FirstName,
					u.LastName,
				),
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserByID)

	router.POST("/users", createUser)
	router.POST("/users/:id", updateUser)

	router.DELETE("/users/:id", deleteUser)

	router.Run(apiUrl)
}
