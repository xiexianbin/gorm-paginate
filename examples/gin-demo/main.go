package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	paginate "go.xiexianbin.cn/gorm-paginate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db, err = gorm.Open(sqlite.Open("demo.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	// Auto Migrate
	db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model

	Name           string
	Age            int
	Balance        int64
	AccountManager string
}

func main() {
	router := gin.Default()

	router.GET("/users", func(c *gin.Context) {
		pagination := paginate.Pagination{}
		tx := db.Scopes(
			paginate.Paginate(
				User{},
				c.Request.URL.Query(),
				&pagination,
				db),
		)

		var users []User
		tx.Find(&users)
		pagination.Items = users
		c.JSON(http.StatusOK, pagination)
	})

	router.POST("/users", func(c *gin.Context) {
		for i := 0; i < 100; i++ {
			user := User{
				Name:           lo.RandomString(5, lo.LowerCaseLettersCharset),
				Age:            i + 10,
				Balance:        int64(i),
				AccountManager: []string{"zhangsi", "lisi", "wangwu"}[i%3],
			}
			db.Create(&user)
		}
		c.JSON(http.StatusAccepted, map[string]string{})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
