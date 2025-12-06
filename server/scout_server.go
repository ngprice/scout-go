package server

import (
	"context"
	pb "scout-go/proto"
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
	game, err := NewGame(int(req.NumPlayers))
	if err != nil {
		return nil, err
	}
	s.games[game.Id] = game

	return &pb.CreateGameResponse{
		GameId: game.Id,
	}, nil
}

func (s *ScoutServer) PlayerAction(ctx context.Context, req *pb.PlayerActionRequest) (*pb.PlayerActionResponse, error) {
	game := s.games[req.GameId]
	err := game.PlayerAction(int(req.PlayerIndex), req.Action, req.Params)
	var msg string
	if err != nil {
		msg = err.Error()
	}
	return &pb.PlayerActionResponse{
		Err:    err != nil,
		ErrMsg: msg,
	}, nil
}

func (s *ScoutServer) GetGameState(ctx context.Context, req *pb.GetGameStateRequest) (*pb.GetGameStateResponse, error) {
	return &pb.GetGameStateResponse{
		Game: s.games[req.GameId].ToProto(),
	}, nil
}

func (s *ScoutServer) GetPlayerState(ctx context.Context, req *pb.GetPlayerStateRequest) (*pb.GetPlayerStateResponse, error) {
	return &pb.GetPlayerStateResponse{
		Player: s.games[req.GameId].Players[req.PlayerIndex].ToProto(),
	}, nil
}
