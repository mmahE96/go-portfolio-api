package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-api-portfolio/models" // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"strconv"

	// used to read the environment variable
	// package used to covert string into int type

	// used to get the params from the route

	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//os.Getenv("POSTGRES_URL")
	// Open the connection
	db, err := sql.Open("postgres", "postgres://ekwwwoarfjjzmq:92e7e65a7e7e439aa49ec5bc45bb1a0a51a49b0705596841003553938655a17d@ec2-34-242-89-204.eu-west-1.compute.amazonaws.com:5432/d6hnola6nvjmb8")
	//postgres://ekwwwoarfjjzmq:92e7e65a7e7e439aa49ec5bc45bb1a0a51a49b0705596841003553938655a17d@ec2-34-242-89-204.eu-west-1.compute.amazonaws.com:5432/d6hnola6nvjmb8

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateUser create a user in the postgres db
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.Article
	var article models.Article

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&article)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := insertArticle(article)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Article created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// insert one user in the DB
func insertArticle(user models.Article) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO articles (title, content, author, date, category, description, slug) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	//psql -h ec2-3-231-112-124.compute-1.amazonaws.com -p 5432-U gmsqtpnkywzlsf dl9bag0ac2qt8
	// the inserted id will store in this id
	var id int64
	id = user.Id

	fmt.Println(id)

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, user.Title, user.Content, user.Author, user.Date, user.Category, user.Description, user.Slug).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// GetArticle will return a single user by its id
func GetArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getArticle function with user id to retrieve a single user
	user, err := getArticle(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(user)
}

func getArticle(id int64) (models.Article, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.Article type
	var article models.Article

	// create the select sql query
	sqlStatement := `SELECT * FROM articles WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&article.Id, &article.Title, &article.Content, &article.Author, &article.Date, &article.Category, &article.Description, &article.Slug)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return article, nil
	case nil:
		return article, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return article, err
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	users, err := getAllArticles()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the users as response
	json.NewEncoder(w).Encode(users)
}

func getAllArticles() ([]models.Article, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var articles []models.Article

	// create the select sql query
	sqlStatement := `SELECT * FROM articles`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var article models.Article

		// unmarshal the row object to user
		err = rows.Scan(&article.Id, &article.Title, &article.Content, &article.Author, &article.Date, &article.Category, &article.Description, &article.Slug)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		articles = append(articles, article)

	}

	// return empty user on error
	return articles, err
}

// DeleteArticle delete user's detail in the postgres db
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("writter,:", w)
	//fmt.Println("request pointer,:", r)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the articleid from the request params, key is "id"
	params := mux.Vars(r)
	//fmt.Println("mux tacka vars (r),:", params)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])
	//fmt.Println("strconv,:", id)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteArticle, convert the int to int64
	deletedRows := deleteArticle(int64(id))

	// format the message string
	msg := fmt.Sprintf("Article removed successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// delete user in the DB
func deleteArticle(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM articles WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
