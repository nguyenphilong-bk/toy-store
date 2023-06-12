package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"toy-store/controllers"
	"toy-store/db"
	"toy-store/forms"
	"toy-store/models"

	"github.com/gin-contrib/gzip"
	uuid "github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// CORSMiddleware ...
// CORS (Cross-Origin Resource Sharing)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// RequestIDMiddleware ...
// Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

var auth = new(controllers.AuthController)

// TokenAuthMiddleware ...
// JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.TokenValid(c)
		c.Next()
	}
}

func main() {
	//Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	//Start the default gin server
	r := gin.Default()

	//Custom form validator
	binding.Validator = new(forms.DefaultValidator)

	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//Start PostgreSQL database
	//Example: db.GetDB() - More info in the models folder
	db.Init()
	err = db.GetDB().AutoMigrate(&models.User{}, models.Cart{}, models.Product{})
	if err != nil {
		fmt.Println("error when migrating models")
	}

	// db.GetDB().AutoMigrate(&models.User{}, &models.Brand{}),

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	// db.InitRedis(1)

	v1 := r.Group("/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)

		v1.POST("/user/login", user.Login)
		v1.POST("/user/register", user.Register)
		v1.GET("/user/logout", user.Logout)
		v1.GET("/user/me", TokenAuthMiddleware(), user.Me)

		/*** START AUTH ***/
		auth := new(controllers.AuthController)

		//Refresh the token when needed to generate new access_token and refresh_token for the user
		v1.POST("/token/refresh", auth.Refresh)

		// Brand APIs
		brand := new(controllers.BrandController)

		v1.POST("/brand", TokenAuthMiddleware(), brand.Create)
		v1.GET("/brands", TokenAuthMiddleware(), brand.All)
		v1.GET("/brand/:id", TokenAuthMiddleware(), brand.One)
		v1.PUT("/brand/:id", TokenAuthMiddleware(), brand.Update)
		v1.DELETE("/brand/:id", TokenAuthMiddleware(), brand.Delete)

		// Category APIs
		category := new(controllers.CategoryController)

		v1.POST("/category", TokenAuthMiddleware(), category.Create)
		v1.GET("/categories", TokenAuthMiddleware(), category.All)
		v1.GET("/category/:id", TokenAuthMiddleware(), category.One)
		v1.PUT("/category/:id", TokenAuthMiddleware(), category.Update)
		v1.DELETE("/category/:id", TokenAuthMiddleware(), category.Delete)

		// Category APIs
		product := new(controllers.ProductController)

		v1.POST("/product", TokenAuthMiddleware(), product.Create)
		v1.GET("/products", TokenAuthMiddleware(), product.All)
		v1.GET("/product/:id", TokenAuthMiddleware(), product.One)
		v1.PUT("/product/:id", TokenAuthMiddleware(), product.Update)
		v1.DELETE("/product/:id", TokenAuthMiddleware(), product.Delete)

		cart := new(controllers.CartController)
		// v1.GET("/cart/me", TokenAuthMiddleware(), cart.MyCart)
		v1.PUT("/cart", TokenAuthMiddleware(), cart.Create)
	}

	r.LoadHTMLGlob("./public/html/*")

	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v0.03",
			"goVersion":             runtime.Version(),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	port := os.Getenv("PORT")

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {

		//Generated using sh generate-certificate.sh
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/myCA.cer",
			KEY:  "./cert/myCA.key",
		}

		r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
	} else {
		r.Run(":" + port)
	}

}
