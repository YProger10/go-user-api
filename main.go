package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    mux := http.NewServeMux()
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok"}`))
    })
    srv := &http.Server{Addr: ":8080", Handler: mux,
        ReadTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second}
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()
    go func() {
        log.Info("listening", "addr", ":8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Error("serve", "err", err); os.Exit(1)
        }
    }()
    <-ctx.Done()
    shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    _ = srv.Shutdown(shutCtx)
    log.Info("shutdown complete")
}
