package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbURI = os.Getenv("MONGO_URI")
var port = os.Getenv("PORT")
var client *mongo.Client
var user = []User{}
var users = AllUser{user}
var token Token

// Note struct
type Note struct {
	ID      string `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string `json:"title,omitempty" bson:"title,omitempty"`
	Content string `json:"content,omitempty" bson:"content,omitempty"`
	Ts      int64  `json:"ts,omitempty" bson:"ts,omitempty"`
	Cdate   string `json:"cDate,omitempty" bson:"cDate,omitempty"`
}

// User struct
type User struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Email string             `json:"email,omitempty" bson:"content,omitempty"`
	MsID  string             `json:"msId,omitempty" bson:"msId,omitempty"`
	Cdate string             `json:"cDate,omitempty" bson:"cDate,omitempty"`
	Notes []int              `json:"notes,omitempty" bson:"notes,omitempty"`
}

// Token struct
type Token struct {
	Refresh     string `json:"refresh_token" bson:"refreshToken,omitempty"`
	Scope       string `json:"scope,omitempty" bson:"scope,omitempty"`
	IDToken     string `json:"id_token,omitempty" bson:"idToken,omitempty"`
	AccessToken string `json:"access_token,omitempty" bson:"accessToken,omitempty"`
}

// AllUser Struct
type AllUser struct {
	Name []User
}

// getUser user
func (s *AllUser) getUser(name string) string {
	var filtered string
	for _, user := range s.Name {
		fmt.Println(user)
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			filtered = name
		}
	}
	return filtered
}

// addUser new user
func (s *AllUser) addUser(user User) []User {
	s.Name = append(s.Name, user)
	return s.Name
}

// func createItem(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Set("Content-Type", "application/json")
// 	var product Product
// 	// query := mux.Vars(req)
// 	product = Product{Desc: query["key"], Price: 1500, Ts: time.Now().UnixNano(), Port: port, Brand: &Brand{Name: "Acetaminofen", Cod: 2000}}
// 	// _ = json.NewDecoder(req.Body).Decode(&product)
// 	collection := client.Database("minimax").Collection("sampleData")
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	result, _ := collection.InsertOne(ctx, product)
// 	userSent := Response{
// 		Date:   123,
// 		Result: &product,
// 	}
// 	fmt.Println(result)
// 	json.NewEncoder(res).Encode(userSent)
// }

// func getAllItem(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Set("Content-Type", "application/json")
// 	var item []Product
// 	collection := client.Database("minimax").Collection("sampleData")
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	cur, err := collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cur.Close(ctx)
// 	for cur.Next(ctx) {
// 		var result Product
// 		err := cur.Decode(&result)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		item = append(item, result)
// 	}
// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	json.NewEncoder(res).Encode(item)
// 	// var product Product
// 	// _ = json.NewDecoder(req.Body).Decode(&product)
// 	// collection := client.Database(`testGo`).Collection(`example`)
// 	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	// result, _ := collection.InsertOne(ctx, product)
// 	// json.NewEncoder(res).Encode(result)
// }
func main() {
	r := gin.Default()

	r.GET("/login", oauth2Log)

	// Serves literal characters
	r.GET("/callback", oauth2Code)

	// Serves Response
	r.GET("/done", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, token)
	})

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		result := users.getUser(id)
		c.PureJSON(http.StatusOK, gin.H{"ok": result})
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func oauth2Log(c *gin.Context) {
	queryStr := url.Values{
		"client_id":     {"b7dbca0c-8bd2-4fd5-bfe1-ed7e7c563f4e"},
		"prompt":        {"select_account"},
		"response_type": {"code"},
		"redirect_uri":  {"http://localhost:8080/callback"},
		"response_mode": {"query"},
		"state":         {"123"},
	}
	apiURL := "https://login.microsoftonline.com"
	resource := "/bf1ce8b5-5d39-4bc5-ad6e-07b3e4d7d67a/oauth2/authorize"
	fmt.Print(queryStr.Encode())
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	u.RawQuery = queryStr.Encode()
	urlStr := u.String()
	c.Redirect(http.StatusMovedPermanently, urlStr)
	c.Abort()
}

func oauth2Code(c *gin.Context) {
	v := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {"b7dbca0c-8bd2-4fd5-bfe1-ed7e7c563f4e"},
		"code":          {c.Query("code")},
		"redirect_uri":  {"http://localhost:8080/callback"},
		"client_secret": {"BI=P]r]KtLW3*9ehjr5JvLbVd3]Fv8N:"},
	}
	resp, _ := http.PostForm("https://login.microsoftonline.com/bf1ce8b5-5d39-4bc5-ad6e-07b3e4d7d67a/oauth2/token", v)
	json.NewDecoder(resp.Body).Decode(&token)

	users.addUser(User{ID: primitive.NewObjectID(), Name: "tester", Email: "eduard.mazo@gmail.com"})
	for i := range users.Name {
		fmt.Println("index:", users.Name[i].Email)
	}
	_, err := c.Cookie("access_token")
	if err != nil {
		c.SetCookie("accessToken", token.AccessToken, 0, "/", "localhost", false, false)
	}

	c.Redirect(http.StatusMovedPermanently, "/done")
	c.Abort()
}
