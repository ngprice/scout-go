package server

import (
	"encoding/json"
	"log"
	pb "scout-go/proto"
)

func (g *Game) ToProto() *pb.Game {
	protoGame := &pb.Game{
		Id:                g.Id,
		ConsecutiveScouts: int32(g.ConsecutiveScouts),
		Round:             int32(g.Round),
		Complete:          g.Complete,
	}

	if g.ActivePlayer != nil {
		protoGame.ActivePlayerIndex = int32(g.ActivePlayer.Index)
	}

	if g.ActiveSetPlayer != nil {
		protoGame.ActiveSetPlayerIndex = int32(g.ActiveSetPlayer.Index)
	}

	for _, card := range g.ActiveSet {
		protoGame.ActiveSet = append(protoGame.ActiveSet, card.ToProto())
	}

	for _, player := range g.Players {
		player_hand := &pb.PlayerHandSize{
			PlayerIndex: int32(player.Index),
			HandSize:    int32(len(player.Hand)),
		}
		protoGame.PlayerHandSize = append(protoGame.PlayerHandSize, player_hand)
	}

	return protoGame
}

func (g *Game) ToJSON() string {
	jg, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		log.Printf("Error marshalling game to JSON: %v", err)
		return ""
	}
	return string(jg)
}

func (p *Player) ToProto() *pb.Player {
	hand := make([]*pb.Card, 0)
	for _, card := range p.Hand {
		hand = append(hand, card.ToProto())
	}

	return &pb.Player{
		Name:            p.Name,
		Index:           int32(p.Index),
		Score:           int32(p.Score),
		Hand:            hand,
		CanReverseHand:  p.CanReverseHand,
		CanScoutAndShow: p.CanScoutAndShow,
	}
}

func (c *Card) ToProto() *pb.Card {
	return &pb.Card{
		Value1: int32(c.Value1),
		Value2: int32(c.Value2),
	}
}

func ToActionSpec(action *pb.Action) *ActionSpec {
	return &ActionSpec{
		ID:             0, // internal use only
		Type:           ActionType(action.ActionType),
		ScoutTakeIndex: int(action.ScoutTakeIndex),
		ScoutPutIndex:  int(action.ScoutPutIndex),
		ShowFirstIndex: int(action.ShowFirstIndex),
		ShowLength:     int(action.ShowLength),
	}
}
