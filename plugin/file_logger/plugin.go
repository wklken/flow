package file_logger

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	plugin_ctx "github.com/wklken/flow/plugin/ctx"
	"go.uber.org/zap"
)

const (
	// version  = "0.1"
	priority = 103
	name     = "file_logger"
)

type Plugin struct {
	config Config
}

// FIXME: use jsonschema to unmarshal the config dynamic

type Config struct{}

func (p *Plugin) Name() string {
	return name
}

func (p *Plugin) Priority() int {
	return priority
}

func (p *Plugin) Init(config string) error {
	fmt.Println("init the request_id plugin", config)
	v := viper.New()
	v.SetConfigType("json")

	// TODO: how to make the default value
	// v.SetDefault("header_name", "X-Request-ID")
	// v.SetDefault("set_in_response", true)

	v.ReadConfig(bytes.NewBuffer([]byte(config)))

	fmt.Println("config: ", v.AllSettings())

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		return err
	}
	fmt.Printf("config: %+v\n", c)
	p.config = c

	return nil
}

func (p *Plugin) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := plugin_ctx.GetPluginVarString(ctx, "request_id", "request_id")

		// TODO: path/filename/rotate
		// TODO: custom fields

		// logger, _ := zap.NewProduction()
		// logger := zap.NewExample()
		// defer logger.Sync()
		// ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		matchedURI := chi.RouteContext(ctx).RoutePattern()
		status := ww.Status()

		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		config.DisableCaller = true
		logger, _ := config.Build()

		fmt.Println("file_logger getting")

		// https://pkg.go.dev/go.uber.org/zap#hdr-Configuring_Zap

		logger.Info("-",
			// Structured context as strongly typed Field values.
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remoteIP", r.RemoteAddr),
			zap.String("proto", r.Proto),
			zap.String("scheme", scheme),
			zap.String("requestID", requestID),
			zap.Int("status", status),
			zap.String("matchedURI", matchedURI),
		)
	}
	return http.HandlerFunc(fn)
}
