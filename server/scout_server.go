package server

import (
	"context"
	"fmt"
	pb "scout-go/proto"
)

const MAX_HAND_SIZE = 20       // practical max
const MAX_ACTIVE_SET_SIZE = 10 // straight 1-10

const (
	ActionScout ActionType = iota
	ActionScoutReverse
	ActionShow
	ActionScoutAndShow
	ActionScoutAndShowReverse
	ActionReverseHand
)

type ActionType int

type ActionSpec struct {
	ID   int
	Type ActionType
	A, B int // meaning depends on type
	C, D int // for scoutAndShow
}

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
		mask[action.ID] = game.IsActionValid(player.Index, action.Type, action.A, action.B, action.C, action.D)
	}

	return &pb.GetValidActionsResponse{Mask: mask}, nil
}

func getAllActions() []ActionSpec {
	actions := make([]ActionSpec, 0)
	scoutActions := make([]ActionSpec, 0)
	scoutActionsReverse := make([]ActionSpec, 0)
	showActions := make([]ActionSpec, 0)

	id := 0

	// 1. scout actions
	for take := 0; take < MAX_ACTIVE_SET_SIZE; take++ {
		for put := 0; put <= MAX_HAND_SIZE; put++ {
			scoutAction := ActionSpec{
				ID: id, Type: ActionScout, A: take, B: put,
			}
			actions = append(actions, scoutAction)
			scoutActions = append(scoutActions, scoutAction)
			id++
		}
	}

	for take := 0; take < MAX_ACTIVE_SET_SIZE; take++ {
		for put := 0; put <= MAX_HAND_SIZE; put++ {
			scoutAction := ActionSpec{
				ID: id, Type: ActionScoutReverse, A: take, B: put,
			}
			actions = append(actions, scoutAction)
			scoutActionsReverse = append(scoutActionsReverse, scoutAction)
			id++
		}
	}

	// 2. show actions
	for start := 0; start < MAX_HAND_SIZE; start++ {
		for length := 1; length <= MAX_HAND_SIZE-start; length++ {
			showAction := ActionSpec{
				ID: id, Type: ActionShow, A: start, B: length,
			}
			actions = append(actions, showAction)
			showActions = append(showActions, showAction)
			id++
		}
	}

	// 3. scoutandshow = Cartesian product
	for _, scout := range scoutActions {
		for _, show := range showActions {
			actions = append(actions, ActionSpec{
				ID:   id,
				Type: ActionScoutAndShow,
				A:    scout.A, B: scout.B,
				C: show.A, D: show.B,
			})
			id++
		}
	}

	// 4. scoutandshowreverse = Cartesian product
	for _, scout := range scoutActionsReverse {
		for _, show := range showActions {
			actions = append(actions, ActionSpec{
				ID:   id,
				Type: ActionScoutAndShowReverse,
				A:    scout.A, B: scout.B,
				C: show.A, D: show.B,
			})
			id++
		}
	}

	// 5. reversehand
	actions = append(actions, ActionSpec{
		ID: id, Type: ActionReverseHand,
	})

	return actions
}
