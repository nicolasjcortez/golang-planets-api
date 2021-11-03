package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"starwars/docs"
	planets "starwars/planets/delivery/http"
	"starwars/planets/repository"
	"starwars/planets/service"
)

func handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, "V.0.0")
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file not loaded \n", err)
		panic(err)
	}
}

// func ApiMiddleware(host, db string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Set("host", host)
// 		c.Set("db", db)
// 		c.Next()
// 	}
// }

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func main() {
	// CLI
	env_arg := flag.String("env", "dev", "Environment eg. prod or qa")
	flag.Parse()
	fmt.Println("\nCONFIG ENVIRONMENT:", *env_arg)

	//CONFIG
	var mongo_host, db, port string
	var protocol string
	mongo_host = viper.GetString(*env_arg + ".MONGO_HOST")
	db = viper.GetString(*env_arg + ".MONGO_DB")
	if *env_arg == "dev" {
		protocol = "http"
		port = viper.GetString(*env_arg + ".PORT")
	} else {
		protocol = "https"
		port = os.Getenv("PORT")
	}

	fmt.Println("MONGO_HOST: ", mongo_host)
	fmt.Println("MONGO_DB: ", db)
	fmt.Println("PORT: ", port)

	// WEB SERVER SETUP
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// r.Use(ApiMiddleware(host, db))
	r.Use(JSONMiddleware())

	// SWAGGER
	docs.SwaggerInfo.Title = "STARWARS PLANETS API " + *env_arg
	docs.SwaggerInfo.Description = "Planets API"
	docs.SwaggerInfo.Version = "V 0.0"
	// docs.SwaggerInfo.Host = "localhost:" + port
	// docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{protocol}

	// ROUTING
	// REST
	r.GET("/", handleVersion)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ph := planets.PlanetHandler{
		PlanetsService: service.PlanetsService{
			Repo: repository.PlanetsRepositoryMongo{
				Host:     mongo_host,
				Database: db,
			},
			Swapi: repository.RemotePlanetsRespositorySwapi{},
		},
	}
	planets.ApplyRoutes(r, ph)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// GRACEFULL SHUTDOWN
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
