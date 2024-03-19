package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/simondrake/mock-azure-ad/internal/handler"
)

func init() {
	build, _ := debug.ReadBuildInfo()

	enableDebug := flag.Bool("debug", false, "whether to enable debug mode")
	flag.Parse()

	opts := &slog.HandlerOptions{AddSource: true}
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
	k := koanf.New(".")

	if err := k.Load(confmap.Provider(map[string]interface{}{
		"server.port":              8080,
		"connection.tenantid":      "mytenantid",
		"connection.certdirectory": "/certs",
		"connection.certname":      "cert.pem",
		"connection.keyname":       "key.pem",
	}, "."), nil); err != nil {
		log.Fatalf("error setting config defaults: %v", err)
	}

	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if err := k.Load(env.Provider("MOCKAZURE_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "MOCKAZURE_")), "_", ".")
	}), nil); err != nil {
		log.Fatalf("error merging environment variables: %v", err)
	}

	serverPort := k.Int("server.port")
	tenantID := k.String("connection.tenantid")
	endpoint := k.String("connection.endpoint")
	certDIR := k.String("connection.certdirectory")
	certName := k.String("connection.certname")
	keyName := k.String("connection.keyname")
	signingKey := k.String("authentication.signingkey")

	slog.Info("Config",
		slog.Group("server",
			slog.Int("port", serverPort),
		),
		slog.Group("connection",
			slog.String("tenantid", tenantID),
			slog.String("endpoint", endpoint),
			slog.String("certdir", certDIR),
			slog.String("certname", certName),
			slog.String("keyname", keyName),
		),
		slog.Group("authentication",
			slog.String("signingkey", signingKey),
		),
	)

	h := handler.New(handler.WithTenantID(tenantID), handler.WithSigningKey(signingKey), handler.WithEndpoint(endpoint))

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("GET /%s/v2.0/.well-known/openid-configuration", tenantID), h.GetOpenID)
	mux.HandleFunc(fmt.Sprintf("POST /%s/oauth2/v2.0/token", tenantID), h.PostToken)

	addr := fmt.Sprintf(":%d", serverPort)
	cert := fmt.Sprintf("%s/%s", certDIR, certName)
	key := fmt.Sprintf("%s/%s", certDIR, keyName)

	slog.Info("Server running", "addr", addr)
	if err := http.ListenAndServeTLS(addr, cert, key, mux); err != nil {
		panic(err)
	}
}
