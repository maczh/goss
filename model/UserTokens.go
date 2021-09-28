package model

type UserTokens struct {
	UserId string      `json:"userId" bson:"userId"`
	Tokens []TokenInfo `json:"tokens" bson:"tokens"`
}
