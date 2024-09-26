// file: main.go

package main

import (
	"fmt"
	"time"

	"go-usip/datasource"
	"go-usip/repositories"
	"go-usip/services"
	"go-usip/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./configs/config.yaml") // the config file path
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	app := iris.New()
	// You got full debug messages, useful when using MVC and you want to make
	// sure that your code is aligned with the Iris' MVC Architecture.
	app.Logger().SetLevel("debug")

	// Load the template files.
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.HandleDir("/public", iris.Dir("./web/public"))

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		if err := ctx.View("shared/error.html"); err != nil {
			ctx.HTML("<h3>%s</h3>", err.Error())
			return
		}
	})

	// ---- Serve our controllers. ----

	// Prepare our repositories and services.
	db, err := datasource.LoadDB()
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}
	userRepo := repositories.NewUserRepository(db)
	fileRepo := repositories.NewFileRepository(db)
	fileCollaRepo := repositories.NewFileCollaboratorRepository(db)

	userService := services.NewUserService(userRepo)
	universerService := services.NewUniverseService()
	fileService := services.NewFileService(fileRepo, fileCollaRepo, universerService)

	sessiondb := redis.New(redis.Config{
		Network:   "tcp",
		Addr:      viper.GetString("redis.addr"),
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Password:  "",
		Database:  "",
		Prefix:    "",
		Driver:    redis.GoRedis(), // redis.Radix() can be used instead.
	})
	// Close connection when control+C/cmd+C
	iris.RegisterOnInterrupt(func() {
		sessiondb.Close()
	})

	sessManager := sessions.New(sessions.Config{
		Cookie:       "_on-premise",
		Expires:      7 * 24 * time.Hour,
		AllowReclaim: true,
	})
	sessManager.UseDatabase(sessiondb)

	// "/user" based mvc application.
	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserController))

	file := mvc.New(app.Party("/file"))
	file.Register(
		fileService,
		sessManager.Start,
	)
	file.Handle(new(controllers.FileController))

	usip := mvc.New(app.Party("/usip"))
	usip.Register(
		userService,
		fileService,
		sessManager.Start,
	)
	usip.Handle(new(controllers.UsipController))

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.Redirect("/file/list", iris.StatusFound)
	})

	// Starts the web server at localhost:8080
	// Enables faster json serialization and more.
	app.Listen(":8080", iris.WithOptimizations)
}
