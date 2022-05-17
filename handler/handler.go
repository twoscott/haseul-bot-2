// Package handler handles events from Discord's API.
package handler

import "github.com/twoscott/haseul-bot-2/router"

// Handler wraps router and handles events from the API, and passes them on
// to the router.
type Handler struct {
	*router.Router
}

// New returns a new instance of Handler.
func New(router *router.Router) *Handler {
	return &Handler{router}
}
