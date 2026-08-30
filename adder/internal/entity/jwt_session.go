package entity

import "github.com/golang-jwt/jwt/v5"

type SessionClaims struct {
	SessionId string `json:"session_id"`
	StationId string `json:"station_id"`
	jwt.RegisteredClaims
}
