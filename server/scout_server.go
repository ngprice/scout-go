package server

import (
	"context"
	pb "scout-ai/proto"

	"github.com/google/uuid"
)

type ScoutServer struct {
	pb.UnimplementedScoutServiceServer
	// games map[string]*Game
}

func NewScoutServer() *ScoutServer {
	return &ScoutServer{}
}

func (s *ScoutServer) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	id := uuid.New().String()

	// game := NewGame(req.NumPlayers)
	// s.games[id] = game

	return &pb.CreateGameResponse{
		GameId: id,
	}, nil
}

func (s *ScoutServer) Step(ctx context.Context, req *pb.StepRequest) (*pb.StepResponse, error) {
	// done := game.ApplyAction(req.Action)

	return &pb.StepResponse{
		Done: false,
	}, nil
}
