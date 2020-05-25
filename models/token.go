package models

import (
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	Token string `json:"token"`
	Valid bool   `json:"valid"`
}

func MarshalToken(token *pb.Token) *Token {
	return &Token{
		Token: token.Token,
		Valid: token.Valid,
	}
}

func UnmarshalToken(token *Token) *pb.Token {
	return &pb.Token{
		Token: token.Token,
		Valid: token.Valid,
	}
}
