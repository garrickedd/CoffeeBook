package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id           int        `json:"id" gorm:"column:id"`
	Username     string     `json:"username" gorm:"column:username"`
	FirstName    string     `json:"first_name" gorm:"column:first_name"`
	LastName     string     `json:"last_name" gorm:"column:last_name"`
	MobileNumber string     `json:"mobile_number" gorm:"column:mobile_number"`
	Email        string     `json:"email" gorm:"column:email"`
	Password     string     `json:"password" gorm:"column:password"`
	Role         int        `json:"role" gorm:"column:role"`
	CreatedAt    *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string { return "users" }

type UserCreation struct {
	Id           int    `json:"-" gorm:"column:id"`
	Username     string `json:"username" gorm:"column:username"`
	FirstName    string `json:"first_name" gorm:"column:first_name"`
	LastName     string `json:"last_name" gorm:"column:last_name"`
	MobileNumber string `json:"mobile_number" gorm:"column:mobile_number"`
	Email        string `json:"email" gorm:"column:email"`
	Password     string `json:"password" gorm:"column:password"`
	Role         int    `json:"role" gorm:"column:role"`
}

func (UserCreation) TableName() string { return User{}.TableName() }

type UserUpdate struct {
	// Id           int    `json:"-" gorm:"column:id"`
	// Username     string `json:"username" gorm:"column:username"`
	FirstName    *string `json:"first_name" gorm:"column:first_name"`
	LastName     *string `json:"last_name" gorm:"column:last_name"`
	MobileNumber *string `json:"mobile_number" gorm:"column:mobile_number"`
	Email        *string `json:"email" gorm:"column:email"`
	Password     *string `json:"password" gorm:"column:password"`
	Role         *int    `json:"role" gorm:"column:role"`
}

func (UserUpdate) TableName() string { return User{}.TableName() }

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit >= 100 {
		p.Limit = 10
	}
}

func main() {
	// GORM connection
	dsn := "admin:password@tcp(127.0.0.1:3306)/coffeebook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// CRUD: Create, Read, Update, Delete
	// POST /v1/users (create new user)
	// GET /v1/users (list user) /v1/users?page=1
	// GET /v1/users/:id (get user details by id)
	// (PUT || PATCH) /v1/users/:id (update user by id)
	// DELETE /v1/users/:id (delete user by id)

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", CreateUser(db))
			users.GET("", ListUser(db))
			users.GET("/:id", GetUser(db))
			users.PATCH("/:id", UpdateUser(db))
			users.DELETE("/:id", DeleteUser(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CreateUser(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data UserCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data.Id,
		})
	}
}

func GetUser(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data User

		id, err := strconv.Atoi(c.Param("id"))
		// fmt.Printf("ID received: %d\n", id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		// data.Id = id
		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func UpdateUser(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data User

		id, err := strconv.Atoi(c.Param("id"))
		// fmt.Printf("ID received: %d\n", id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func DeleteUser(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Table(User{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
			"role": -1,
		}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func ListUser(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var result []User

		// filter if role -1
		db = db.Where("role <> ?", -1)

		if err := db.Table(User{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		if err := db.Order("id").
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&result).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   result,
			"paging": paging,
		})
	}
}
