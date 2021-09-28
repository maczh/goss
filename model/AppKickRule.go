package model

type AppKickRule struct {
	AppId          string          `json:"appId" bson:"appId"`
	KickRule       int             `json:"kickRule" bson:"kickRule"`             //互踢规则类型 1-同端互踢 2-完全互踢，单端在线 3-指定端互踢 4-指定端不互踢
	TermTypes      []int           `json:"termTypes" bson:"termTypes"`           //互踢或互不踢的终端类型列表
	TermTypeGroups []TermTypeGroup `json:"termTypeGroups" bson:"termTypeGroups"` //互踢组列表
}
