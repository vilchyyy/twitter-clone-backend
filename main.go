package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type comment struct {
	TID     string   `json:"tid"`
	ID      string   `json:"id"`
	Image   string   `json:"image"`
	AccName string   `json:"accName"`
	AtName  string   `json:"atName"`
	Text    string   `json:"text"`
	Likes   []string `json:"likes"`
}

type tweet struct {
	ID       string    `json:"id"`
	Image    string    `json:"image"`
	AccName  string    `json:"accName"`
	AtName   string    `json:"atName"`
	Text     string    `json:"text"`
	Likes    []string  `json:"likes"`
	Shares   int       `json:"shares"`
	Comments []comment `json:"Comments"`
}

type like struct {
	ID      string `json:"tid"`
	AccName string `json:"accName"`
}

type commentLike struct {
	TID     string `json:"tid"`
	ID      string `json:"id"`
	AccName string `json:"accName"`
}

var tweets = []tweet{}

func index(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getTweets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tweets)
}

func postTweets(c *gin.Context) {
	var newTweet tweet

	if err := c.BindJSON(&newTweet); err != nil {
		return
	}

	tweets = append([]tweet{newTweet}, tweets...)
	c.IndentedJSON(http.StatusCreated, newTweet)
}

func postComments(c *gin.Context) {
	var newComment comment

	if err := c.BindJSON(&newComment); err != nil {
		return
	}
	for i := range tweets {
		a := &tweets[i]
		if a.ID == newComment.TID {
			a.Comments = append([]comment{newComment}, a.Comments...)
			c.IndentedJSON(http.StatusOK, newComment)
		}
	}

}

func likePost(c *gin.Context) {
	var newLike like

	if err := c.BindJSON(&newLike); err != nil {
		return
	}
	for i := range tweets {
		a := &tweets[i]
		if a.ID == newLike.ID {
			ind := index(a.Likes, newLike.AccName)
			if ind == -1 {
				a.Likes = append(a.Likes, newLike.AccName)
			} else {
				a.Likes = remove(a.Likes, ind)
			}
			c.IndentedJSON(http.StatusOK, a)
		}
	}
}

func likeComment(c *gin.Context) {
	var newCommentLike commentLike

	if err := c.BindJSON(&newCommentLike); err != nil {
		return
	}
	for i := range tweets {
		a := &tweets[i]
		if a.ID == newCommentLike.TID {
			for j := range a.Comments {
				b := &a.Comments[j]
				if b.ID == newCommentLike.ID {
					ind := index(b.Likes, newCommentLike.AccName)
					if ind == -1 {
						b.Likes = append(b.Likes, newCommentLike.AccName)
					} else {
						b.Likes = remove(b.Likes, ind)
					}
					c.IndentedJSON(http.StatusOK, b)
				}
			}
		}
	}
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.PUT("/tweets", postComments)
	router.PUT("/tweets/like", likePost)
	router.PUT("/tweets/comments/like", likeComment)
	router.GET("/tweets", getTweets)
	router.POST("/tweets", postTweets)
	router.Run(fmt.Sprintf(":%s", "127.0.0.1:8080"))
}
