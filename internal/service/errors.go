package service

import "errors"

var (
	ErrForbidden     = errors.New("tenant_mismatch")
	ErrMissingTenant = errors.New("missing_tenant")
	ErrNotFound      = errors.New("not_found")
	ErrInvalidInput  = errors.New("invalid_input")
)
