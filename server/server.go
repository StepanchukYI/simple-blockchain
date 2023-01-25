package server

import "context"

type Server interface {
	Status() error
	Serve() error
	Shutdown(ctx context.Context) error
}
