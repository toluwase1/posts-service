package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Post struct {
	Id          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Comments    []Comment `json:"comments" gorm:"-" default:"[]"`
}

type Comment struct {
	Id     uint   `json:"id"`
	PostId uint   `json:"post_id"`
	Text   string `json:"text"`
}

func main() {
	db, err := gorm.Open(mysql.Open("root:toluwase@tcp(127.0.0.1:3306)/posts_ms"), &gorm.Config{})
	if err != nil {
		panic("database error")
	}
	db.AutoMigrate(Post{})

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/api/posts", func(c *fiber.Ctx) error {
		var posts []Post

		for i, post := range posts {
			response, err := http.Get(fmt.Sprintf("http://localhost:8001/api/posts/%d/comments", post.Id))
			if err != nil {
				log.Println(err)
				return err
			}
			var comment []Comment

			json.NewDecoder(response.Body).Decode(&comment)
			posts[i].Comments = comment
		}
		db.Find(&posts)
		return c.JSON(posts)
	})

	app.Post("/api/posts", func(c *fiber.Ctx) error {
		var post Post

		if err := c.BodyParser(&post); err != nil {
			log.Println(err)
			return err
		}
		db.Create(&post)
		return c.JSON(post)
	})

	app.Listen(":8000")

}
