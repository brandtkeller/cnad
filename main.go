package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	remoteURL = "https://test" // endpoint to try first
	addr      = "0.0.0.0:8080"
)

var (
	logger     = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	serverHTML string // cached content to serve
)

func loadContent() string {
	client := &http.Client{Timeout: 3 * time.Second}

	// Try remote content
	if resp, err := client.Get(remoteURL); err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			logger.Info("loaded remote content", "url", remoteURL)
			return string(body)
		}
		logger.Error("failed reading remote response", "error", err)
	} else if err != nil {
		logger.Warn("remote content fetch failed", "url", remoteURL, "error", err)
	} else {
		logger.Warn("remote content returned non-OK status", "url", remoteURL, "status", resp.StatusCode)
	}

	// Support your local airgaps
	if file, err := os.ReadFile("local.html"); err != nil {
		logger.Error("failed to read local file")
	} else {
		logger.Info("local file content found")
		return string(file)
	}

	// Default failure content
	logger.Info("using default failure HTML")
	return `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Offline Mode</title>
			<style>
				body { font-family: sans-serif; background-color: #0b0c10; color: #c5c6c7; text-align: center; padding: 5em; }
				h1 { color: #66fcf1; }
				p { margin-top: 1em; font-size: 1.2em; }
			</style>
		</head>
		<body>
			<h1>⚠️ Failed to retrieve website content</h1>
			<p>You appear to be offline or the remote site is unavailable.</p>
			<p>This is the fallback user experience in a disconnected environment.</p>
		</body>
		</html>`
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger.Info("serving request", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, serverHTML)
}

func main() {
	// Preload content on startup
	serverHTML = loadContent()

	// Start webserver
	http.HandleFunc("/", handler)
	logger.Info("starting server", "listen_address", addr, "url", "http://"+addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("server stopped unexpectedly", "error", err)
	}
}
