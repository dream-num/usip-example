// file: main.go

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
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

const defaultPort = 8090

func resolvePort() (port int, source string, err error) {
	if value := strings.TrimSpace(os.Getenv("PORT")); value != "" {
		port, err = strconv.Atoi(value)
		if err != nil || port <= 0 || port > 65535 {
			return 0, "", fmt.Errorf("invalid PORT: %q", value)
		}
		return port, "env:PORT", nil
	}

	if value := strings.TrimSpace(viper.GetString("server.port")); value != "" {
		port, err = strconv.Atoi(value)
		if err != nil || port <= 0 || port > 65535 {
			return 0, "", fmt.Errorf("invalid server.port: %q", value)
		}
		return port, "config:server.port", nil
	}

	return defaultPort, "default", nil
}

func resolveRedisEnabled() (enabled bool, source string, err error) {
	if value := strings.TrimSpace(os.Getenv("REDIS_ENABLED")); value != "" {
		enabled, err = strconv.ParseBool(value)
		if err != nil {
			return false, "", fmt.Errorf("invalid REDIS_ENABLED: %q", value)
		}
		return enabled, "env:REDIS_ENABLED", nil
	}

	return viper.GetBool("redis.enabled"), "config:redis.enabled", nil
}

func newUniverserProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(strings.TrimSpace(target))
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
	}
	return proxy, nil
}

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
	app.Use(func(ctx iris.Context) {
		path := ctx.Path()
		if strings.HasPrefix(path, "/sheet") || path == "/files" || path == "/login" || path == "/register" {
			ctx.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			ctx.Header("Pragma", "no-cache")
			ctx.Header("Expires", "0")
		}
		ctx.Next()
	})

	app.HandleDir("/public", iris.Dir("./web/public"))
	app.HandleDir("/sheet", iris.Dir("./web/public/sheet-host"))
	app.Get("/sheet", func(ctx iris.Context) {
		ctx.ServeFile("./web/public/sheet-host/index.html")
	})
	app.Get("/sheet/", func(ctx iris.Context) {
		ctx.ServeFile("./web/public/sheet-host/index.html")
	})
	app.Get("/login", func(ctx iris.Context) {
		ctx.ServeFile("./web/public/sheet-host/index.html")
	})
	app.Get("/register", func(ctx iris.Context) {
		ctx.ServeFile("./web/public/sheet-host/index.html")
	})
	app.Get("/files", func(ctx iris.Context) {
		ctx.ServeFile("./web/public/sheet-host/index.html")
	})

	universerProxy, err := newUniverserProxy(viper.GetString("universer.host"))
	if err != nil {
		app.Logger().Fatalf("invalid universer.host for proxy: %v", err)
		return
	}
	app.Any("/universer-api", iris.FromStd(universerProxy))
	app.Any("/universer-api/{path:path}", iris.FromStd(universerProxy))

	app.OnAnyErrorCode(func(ctx iris.Context) {
		message := ctx.Values().GetStringDefault("message", "The page you're looking for doesn't exist")
		ctx.StopWithJSON(ctx.GetStatusCode(), iris.Map{
			"error":   message,
			"status":  ctx.GetStatusCode(),
			"path":    ctx.Path(),
			"method":  ctx.Method(),
			"request": ctx.GetID(),
		})
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

	avatarService := services.NewAvatarService()
	userService := services.NewUserService(userRepo, avatarService)
	universerService := services.NewUniverseService()
	fileService := services.NewFileService(fileRepo, fileCollaRepo, universerService)

	sessManager := sessions.New(sessions.Config{
		Cookie:                      "_on-premise",
		Expires:                     7 * 24 * time.Hour,
		AllowReclaim:                true,
		DisableSubdomainPersistence: true,
	})

	redisEnabled, redisSource, err := resolveRedisEnabled()
	if err != nil {
		app.Logger().Fatal(err.Error())
		return
	}
	app.Logger().Infof("redis enabled resolved from %s: %t", redisSource, redisEnabled)

	if redisEnabled {
		redisAddr := strings.TrimSpace(viper.GetString("redis.addr"))
		if redisAddr == "" {
			app.Logger().Fatal("redis is enabled but redis.addr is empty")
			return
		}

		sessiondb := redis.New(redis.Config{
			Network:   "tcp",
			Addr:      redisAddr,
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
		sessManager.UseDatabase(sessiondb)
	} else {
		app.Logger().Warn("redis disabled; using in-memory session storage")
	}

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

	authAPI := mvc.New(app.Party("/api/auth"))
	authAPI.Register(
		userService,
		sessManager.Start,
	)
	authAPI.Handle(new(controllers.AuthAPIController))

	filesAPI := mvc.New(app.Party("/api/files"))
	filesAPI.Register(
		fileService,
		sessManager.Start,
	)
	filesAPI.Handle(new(controllers.FilesAPIController))

	usip := mvc.New(app.Party("/usip"))
	usip.Register(
		userService,
		fileService,
		sessManager.Start,
	)
	usip.Handle(new(controllers.UsipController))

	cors := mvc.New(app.Party("/cors"))
	cors.Register(sessManager.Start)
	cors.Handle(new(controllers.CorsController))

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.Redirect("/files", iris.StatusFound)
	})

	port, portSource, err := resolvePort()
	if err != nil {
		app.Logger().Fatal(err.Error())
		return
	}
	app.Logger().Infof("listen port resolved from %s: %d", portSource, port)

	// Starts the web server on configured port
	// Enables faster json serialization and more.
	app.Listen(fmt.Sprintf(":%d", port), iris.WithOptimizations)
}
