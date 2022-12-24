package gapi

import (
	"fmt"

	db "github.com/KhanbalaRashidov/SimpleBank/db/sqlc"
	"github.com/KhanbalaRashidov/SimpleBank/pb"
	"github.com/KhanbalaRashidov/SimpleBank/token"
	"github.com/KhanbalaRashidov/SimpleBank/util"
	"github.com/KhanbalaRashidov/SimpleBank/worker"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
