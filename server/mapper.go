package server

import (
	"encoding/json"
	"log"
	pb "scout-ai/proto"
)

func (g *Game) ToProto() *pb.Game {
	protoGame := &pb.Game{}

	for _, player := range g.Players {
		protoGame.Players = append(protoGame.Players, player.ToProto())
	}

	for _, card := range g.ActiveSet {
		protoGame.ActiveSet = append(protoGame.ActiveSet, card.ToProto())
	}

	protoGame.ActivePlayer = g.ActivePlayer.ToProto()

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
		Name:              p.Name,
		Index:             int32(p.Index),
		Score:             int32(p.Score),
		FlipHandAvail:     p.FlipHandAvail,
		ScoutAndShowAvail: p.ScoutAndShowAvail,
		Hand:              hand,
	}
}

func (c *Card) ToProto() *pb.Card {
	return &pb.Card{
		Value1: int32(c.Value1),
		Value2: int32(c.Value2),
	}
}
