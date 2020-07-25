package models

import "time"

type Measurement struct {
	Id          int64     `json:"id"`
	TimeStamp   time.Time `json:"timestamp"`
	Pressure    float64   `json:"pressure"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
}

type MeasuringNode struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	IsPublic   bool    `json:"is_public"`
	IsOutdoors bool    `json:"is_outdoors"`
}

type NodeAuthToken struct {
	Id           int64
	TokenHash    []byte
	CreationTime time.Time
}

type User struct {
	Id               int64     `json:"id"`
	LastLogin        time.Time `json:"last-login"`
	CreationTime     time.Time `json:"creation-time"`
	Email            string    `json:"email"`
	Username         string    `json:"username"`
	IsEnabled        bool      `json:"enabled"`
	EnableSecretHash []byte    `json:"-"`
	PasswordHash     []byte    `json:"-"`
}

type Invitation struct {
	Id           int64     `json:"id"`
	Email        string    `json:"email"`
	CreationTime time.Time `json:"creation-time"`
}
