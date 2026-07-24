package bootstrap

import (
	"log"
	"net/http"
	"os"
	"time"

	"myapp/config"
)

// NewLogger builds the application-wide logger. In debug mode it writes to
// stdout only; production logging (to storage/logs) is extended in a later
// phase.
func NewLogger(cfg *config.Config) *log.Logger {
	prefix := "[" + cfg.App.Name + "] "
	return log.New(os.Stdout, prefix, log.Ldate|log.Ltime)
}

// NewRouter builds the top-level HTTP handler. Phase 1 uses the standard
// library ServeMux directly; it is replaced by the framework router in a
// later phase once routes/middleware are introduced.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body := `{"status":"ok","time":"` + time.Now().UTC().Format(time.RFC3339) + `"}`
	_, _ = w.Write([]byte(body))
}
