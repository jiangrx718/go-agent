package gins

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"go-agent/gopkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes()
}

type HttpServer struct {
	http.Server
	router   *gin.Engine
	handlers []Handler
}

func NewHttpServer(listen string) *HttpServer {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 设置信任的代理 IP（CIDR 格式）
	if err := r.SetTrustedProxies([]string{
		"10.0.0.0/8",     // 内网 IP 段
		"172.16.0.0/12",  // 内网 IP 段
		"192.168.0.0/16", // 内网 IP 段
	}); err != nil {
		panic(err)
	}

	// 加载静态页面
	r.StaticFile("/s3", "./html/s3.html")
	r.StaticFile("/chat", "./html/chat.html")

	r.Use(gin.Recovery(), RequestID())
	if utils.Debug() {
		gin.SetMode(gin.DebugMode)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	srv := &HttpServer{
		router: r,
		Server: http.Server{
			Addr:    listen,
			Handler: r,
		},
		handlers: []Handler{},
	}

	return srv
}

func (s *HttpServer) RegisterHandler(funcs ...func(*gin.Engine) Handler) {
	for _, fun := range funcs {
		s.handlers = append(s.handlers, fun(s.router))
	}

	for _, handler := range s.handlers {
		handler.RegisterRoutes()
	}
}

func (s *HttpServer) GracefulStart(ctx context.Context) {
	go func() {
		log.Printf("Server listen on %s\n", s.Addr)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	c, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.Shutdown(c); err != nil {
		log.Printf("Server Shutdown: %s\n", err)
	}
	log.Println("Server exiting")
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := strings.ReplaceAll(c.GetHeader("x-request-id"), "-", "")
		if id == "" {
			id = utils.GenUUIDWithoutUnderline()
		}
		c.Set("x-request-id", id)

		c.Next()

		c.Header("x-request_id", id)
	}
}
