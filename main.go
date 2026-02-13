package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
)

const (
	DEFAULT_PORT = ":9737"
)

func main() {
	if debug := os.Getenv("DEBUG"); debug == "1" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	port := DEFAULT_PORT
	if portByEnv := os.Getenv("SSE_PORT"); portByEnv != "" {
		port = fmt.Sprintf(":%s", portByEnv)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors.Default())

	router.GET("/sse", func(ctx *gin.Context) {
		sse.Encode(ctx.Writer, sse.Event{Event: "connection", Data: "open"})
		ctx.Writer.Flush()
		defer func() {
			slog.Debug("Closed")
			sse.Encode(ctx.Writer, sse.Event{Event: "connection", Data: "close"})
			ctx.Writer.Flush()
		}()

		delayStr := ctx.Query("delay")
		delayInt, err := strconv.Atoi(delayStr)
		if err != nil {
			if delayStr != "" {
				slog.Error("delay conversion fail", "err", err)
			}
			sse.Encode(ctx.Writer, sse.Event{Event: "err", Data: "Minimum delay is 1ms (/sse?delay=1)"})
			return
		}

		for {
			select {
			case <-ctx.Writer.CloseNotify():
				return
			default:
			}

			slog.Debug("Sent")
			sse.Encode(ctx.Writer, sse.Event{
				Event: "data",
				Data:  map[string]any{"name": gofakeit.Name()},
			})
			ctx.Writer.Flush()
			time.Sleep(time.Duration(delayInt * int(time.Millisecond)))
		}
	})

	slog.Info("SSE will active with", "port", port)
	if err := router.Run(port); err != nil {
		slog.Error("init sse failed", "err", err)
		return
	}
}
