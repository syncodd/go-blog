package main

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

var db *sql.DB

type Blog struct {
	ID      int64
	Title   string
	Content string
}

func connectMySQL() {

	cfg := mysql.Config{
		User:                 "root",     //os.Getenv("DBUSER"),
		Passwd:               "12345678", //os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "blogs",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

}

func getBlogs() ([]Blog, error) {
	// An albums slice to hold data from returned rows.
	var blogs []Blog

	rows, err := db.Query("SELECT * FROM blog ORDER BY id DESC")
	if err != nil {
		return nil, fmt.Errorf("Blog: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var blg Blog
		if err := rows.Scan(&blg.ID, &blg.Title, &blg.Content); err != nil {
			return nil, fmt.Errorf("Blog: %v", err)
		}
		blogs = append(blogs, blg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Blog: %v", err)
	}
	return blogs, nil
}

func sendBlog(paylaod struct {
	Title   string `form:"title_text"`
	Content string `form:"content_text"`
}) (int64, error) {

	result, err := db.Exec("INSERT INTO blog (title, content) VALUES (?, ?)", paylaod.Title, paylaod.Content)
	if err != nil {
		return 0, fmt.Errorf("addBlog: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addBlog: %v", err)
	}
	return id, nil
}

func delBlog(id string) error {

	_, err := db.Exec("DELETE FROM blog WHERE id = (?)", id)
	if err != nil {
		return fmt.Errorf("addBlog: %v", err)
	}

	return nil
}

func main() {

	connectMySQL()

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {

		blogs, err := getBlogs()
		if err != nil {
			log.Fatal(err)
		}

		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Title",
			"Blogs": blogs,
		})
	})

	app.Get("/blogs", func(c *fiber.Ctx) error {

		blogs, err := getBlogs()
		if err != nil {
			log.Fatal(err)
		}

		// Render index template
		return c.Render("blogs", fiber.Map{
			"Title": "Title",
			"Blogs": blogs,
		})
	})

	app.Post("/blogs", func(c *fiber.Ctx) error {

		payload := struct {
			Title   string `form:"title_text"`
			Content string `form:"content_text"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		sendBlog(payload)

		blogs, err := getBlogs()
		if err != nil {
			log.Fatal(err)
		}

		// Render index template
		return c.Render("blogs", fiber.Map{
			"Title": "Title",
			"Blogs": blogs,
		})
	})

	app.Get("/blogs/delete/:id", func(c *fiber.Ctx) error {

		blog_id := c.Params("id")

		err := delBlog(blog_id)
		if err != nil {
			log.Fatal(err)
		}

		blogs, err := getBlogs()
		if err != nil {
			log.Fatal(err)
		}

		// Render index template
		return c.Render("blogs", fiber.Map{
			"Title": "Title",
			"Blogs": blogs,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
