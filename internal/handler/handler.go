package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type HandlerOptions func(*Handler)

func WithSigningKey(sk string) HandlerOptions {
	return func(h *Handler) {
		h.signingKey = sk
	}
}

func WithTenantID(tid string) HandlerOptions {
	return func(h *Handler) {
		h.tenantID = tid
	}
}

func WithEndpoint(e string) HandlerOptions {
	return func(h *Handler) {
		h.endpoint = e
	}
}

type Handler struct {
	signingKey string
	tenantID   string
	endpoint   string
}

func New(opts ...HandlerOptions) *Handler {
	h := &Handler{}

	for _, o := range opts {
		o(h)
	}

	return h
}

func (h *Handler) GetOpenID(w http.ResponseWriter, r *http.Request) {
	out := map[string]interface{}{
		"issuer":                                h.endpoint,
		"authorization_endpoint":                fmt.Sprintf("%s/%s/oauth2/v2.0/authorize", h.endpoint, h.tenantID),
		"token_endpoint":                        fmt.Sprintf("%s/%s/oauth2/v2.0/token", h.endpoint, h.tenantID),
		"token_endpoint_auth_methods_supported": []string{"client_secret_post", "private_key_jwt"},
		"jwks_uri":                              fmt.Sprintf("%s/%s/discovery/v2.0/keys", h.endpoint, h.tenantID),
		"userinfo_endpoint":                     "https://graph.microsoft.com/oidc/userinfo",
		"subject_types_supported":               []string{"pairwise"},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "failed to encode response: %v", http.StatusInternalServerError)
	}
}

func (h *Handler) PostToken(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().AddDate(1, 0, 0).Unix(),
		"iss": "https://sts.windows-ppe.net/ab1f708d-50f6-404c-a006-d71b2ac7a606/",
		"aud": "https://storage.azure.com",
		"scp": "user_impersonation",
	})
	sigtkn, err := tkn.SignedString([]byte(h.signingKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := map[string]string{
		"token_type":     "Bearer",
		"expires_in":     "3599",
		"ext_expires_in": "3599",
		"access_token":   sigtkn,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "failed to encode response: %v", http.StatusInternalServerError)
	}
}
