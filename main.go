package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/simondrake/mock-azure-ad/internal/handler"
)

func init() {
	build, _ := debug.ReadBuildInfo()

	enableDebug := flag.Bool("debug", false, "whether to enable debug mode")
	flag.Parse()

	opts := &slog.HandlerOptions{}
	if *enableDebug {
		opts.Level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)).With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", build.GoVersion),
		),
	))
}

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	tenantID := os.Getenv("CONNECTION_TENANT_ID")
	endpoint := os.Getenv("CONNECTION_ENDPOINT")
	signingKey := os.Getenv("AUTHENTICATION_SIGNING_KEY")

	if serverPort == "" {
		serverPort = "8080"
	}

	if tenantID == "" {
		tenantID = "mytenantid"
	}

	slog.Info("Config",
		slog.Group("server", slog.String("port", serverPort)),
		slog.Group("connection", slog.String("tenantid", tenantID), slog.String("endpoint", endpoint)),
		slog.Group("authentication", slog.String("signingkey", signingKey)),
	)

	opts := []handler.HandlerOptions{handler.WithTenantID(tenantID)}
	if signingKey != "" {
		opts = append(opts, handler.WithSigningKey(signingKey))
	}

	if endpoint != "" {
		opts = append(opts, handler.WithEndpoint(endpoint))
	}

	h := handler.New(opts...)

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("GET /%s/v2.0/.well-known/openid-configuration", tenantID), h.GetOpenID)
	mux.HandleFunc(fmt.Sprintf("POST /%s/oauth2/v2.0/token", tenantID), h.PostToken)

	slog.Info("Server running", "addr", fmt.Sprintf(":%s", serverPort))
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%s", serverPort), "/certs/cert.pem", "/certs/key.pem", mux); err != nil {
		panic(err)
	}
}
