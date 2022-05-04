package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var App *Services

type Services struct {
	R    *gin.Engine
	C    *GinConfig
	D    *gorm.DB
	MODE string
}

func NewServices(R *gin.Engine, C *GinConfig, D *gorm.DB, MODE string) *Services {
	return &Services{
		R:    R,
		C:    C,
		D:    D,
		MODE: MODE,
	}
}

func (s *Services) GetServerMode() string {
	return s.MODE
}
