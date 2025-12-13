package server

import (
	"context"
	"fmt"
	pb "scout-go/proto"
)

const MAX_HAND_SIZE = 20       // practical max
const MAX_ACTIVE_SET_SIZE = 10 // straight 1-10

type ScoutServer struct {
	pb.UnimplementedScoutServiceServer
	Games      map[string]*Game
	AllActions []ActionSpec
}

func NewScoutServer() *ScoutServer {
	return &ScoutServer{
		Games:      make(map[string]*Game),
		AllActions: getAllActions(),
	}
}

func (s *ScoutServer) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	game, err := NewGame(int(req.NumPlayers))
	if err != nil {
		return nil, err
	}
	s.Games[game.Id] = game

	return &pb.CreateGameResponse{
		GameId: game.Id,
	}, nil
}

func (s *ScoutServer) PlayerAction(ctx context.Context, req *pb.PlayerActionRequest) (*pb.PlayerActionResponse, error) {
	game := s.Games[req.GameId]
	if game == nil {
		return nil, fmt.Errorf("invalid game_id")
	}

	err := game.PlayerAction(int(req.PlayerIndex), ToActionSpec(req.Action))

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
		Game: s.Games[req.GameId].ToProto(),
	}, nil
}

func (s *ScoutServer) GetPlayerState(ctx context.Context, req *pb.GetPlayerStateRequest) (*pb.GetPlayerStateResponse, error) {
	return &pb.GetPlayerStateResponse{
		Player: s.Games[req.GameId].Players[req.PlayerIndex].ToProto(),
	}, nil
}

func (s *ScoutServer) GetValidActions(ctx context.Context, req *pb.GetValidActionsRequest) (*pb.GetValidActionsResponse, error) {
	game := s.Games[req.GameId]
	player := game.Players[req.PlayerIndex]

	mask := make([]bool, len(s.AllActions))

	for _, action := range s.AllActions {
		mask[action.ID] = game.IsActionValid(player.Index, &action)
	}

	return &pb.GetValidActionsResponse{Mask: mask}, nil
}
