package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/ichigozero/gtdzero/controllers"
	"github.com/ichigozero/gtdzero/db/gorm"
	"github.com/ichigozero/gtdzero/libs/auth"
	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/routers"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	libgorm "gorm.io/gorm"
)

func main() {
	fs := flag.NewFlagSet("gtdzero", flag.ExitOnError)
	var (
		httpAddr = fs.String(
			"grpc.addr",
			getEnv("HTTP_ADDR", ":8080"),
			"gRPC listen address",
		)
		redisURL = fs.String(
			"redis.url",
			getEnv("REDIS_URL", ":6379"),
			"Redis URL",
		)
		databaseURL = fs.String(
			"database.url",
			getEnv("DATABASE_URL", ""),
			"Database URL",
		)
	)

	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	client := redis.NewClient(&redis.Options{
		Addr: *redisURL,
	})

	var db *libgorm.DB
	var err error
	{
		if *databaseURL != "" {
			db, err = libgorm.Open(postgres.Open(*databaseURL), &libgorm.Config{})
		} else {
			db, err = libgorm.Open(sqlite.Open("gorm.db"), &libgorm.Config{})
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	db.AutoMigrate(models.User{}, models.Task{})

	userDB := gorm.NewUserDB(db)
	taskDB := gorm.NewTaskDB(db)
	tokenizer := auth.NewTokenizer()
	authClient := auth.NewAuthClient(client)

	r := gin.Default()
	ac := controllers.NewAuthController(userDB, tokenizer, authClient)
	tc := controllers.NewTaskController(taskDB, authClient)

	routers.SetAuthRoutes(r, ac)
	routers.SetTaskRoutes(r, tc)

	r.Run(*httpAddr)
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
