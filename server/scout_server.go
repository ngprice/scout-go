package server

import (
	"context"
	"fmt"
	"sync"

	pb "scout-go/proto"
)

type ScoutServer struct {
	pb.UnimplementedScoutServiceServer

	mu         sync.RWMutex
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

	s.mu.Lock()
	s.Games[game.Id] = game
	s.mu.Unlock()

	return &pb.CreateGameResponse{GameId: game.Id}, nil
}

func (s *ScoutServer) PlayerAction(ctx context.Context, req *pb.PlayerActionRequest) (*pb.PlayerActionResponse, error) {
	// Read-lock just to find the game pointer.
	s.mu.RLock()
	game := s.Games[req.GameId]
	s.mu.RUnlock()

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
	s.mu.RLock()
	game := s.Games[req.GameId]
	s.mu.RUnlock()

	if game == nil {
		return nil, fmt.Errorf("invalid game_id")
	}

	// game.ToProto should be safe (either it locks internally or ToProto snapshots)
	return &pb.GetGameStateResponse{Game: game.ToProto()}, nil
}

func (s *ScoutServer) GetPlayerState(ctx context.Context, req *pb.GetPlayerStateRequest) (*pb.GetPlayerStateResponse, error) {
	s.mu.RLock()
	game := s.Games[req.GameId]
	s.mu.RUnlock()

	if game == nil {
		return nil, fmt.Errorf("invalid game_id")
	}
	if int(req.PlayerIndex) < 0 || int(req.PlayerIndex) >= len(game.Players) {
		return nil, fmt.Errorf("invalid player_index")
	}

	// game.Players[...] access should be safe (again: game-level lock ideally)
	return &pb.GetPlayerStateResponse{
		Player: game.Players[req.PlayerIndex].ToProto(),
	}, nil
}

func (s *ScoutServer) GetValidActions(ctx context.Context, req *pb.GetValidActionsRequest) (*pb.GetValidActionsResponse, error) {
	s.mu.RLock()
	game := s.Games[req.GameId]
	s.mu.RUnlock()

	if game == nil {
		return nil, fmt.Errorf("invalid game_id")
	}
	if int(req.PlayerIndex) < 0 || int(req.PlayerIndex) >= len(game.Players) {
		return nil, fmt.Errorf("invalid player_index")
	}

	player := game.Players[req.PlayerIndex]
	mask := make([]bool, len(s.AllActions))

	for _, action := range s.AllActions {
		mask[action.ID] = game.IsActionValid(player.Index, &action)
	}

	return &pb.GetValidActionsResponse{Mask: mask}, nil
}
