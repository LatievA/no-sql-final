// Package docs bicycle-store API.
//
// Documentation for Bicycle Store API
//
//	Schemes: http
//	BasePath: /api/v1
//	Version: 1.0.0
//	Host: localhost:8080
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- BearerAuth:
//
//	SecurityDefinitions:
//	BearerAuth:
//	  type: apiKey
//	  in: header
//	  name: Authorization
//
// swagger:meta
package docs

import "bicycle-store/internal/models"

// swagger:response errorResponse
type errorResponse struct {
	// in:body
	Body models.APIResponse
}

// swagger:response successResponse
type successResponse struct {
	// in:body
	Body models.APIResponse
}
