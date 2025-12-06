package server

import (
	"context"
	pb "scout-go/proto"

	"github.com/google/uuid"
)

type ScoutServer struct {
	pb.UnimplementedScoutServiceServer
	games map[string]*Game
}

func NewScoutServer() *ScoutServer {
	return &ScoutServer{
		games: make(map[string]*Game),
	}
}

func (s *ScoutServer) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	id := uuid.New().String()

	game := NewGame(int(req.NumPlayers))
	s.games[id] = game

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

func (s *ScoutServer) GetGameState(ctx context.Context, req *pb.GetGameStateRequest) (*pb.GetGameStateResponse, error) {
	response := &pb.GetGameStateResponse{}

	response.Game = s.games[req.GameId].ToProto()

	return response, nil
}
