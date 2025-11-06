package handler

import (
    "net/http"
)

// PreflightHandler returns 200 for CORS preflight requests.
func PreflightHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Headers are set by global CORS middleware; just acknowledge.
        w.WriteHeader(http.StatusOK)
    }
}