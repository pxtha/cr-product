package route

import (
	"github.com/caarlos0/env"
	"gitlab.com/goxp/cloud0/service"
)

type extraSetting struct {
	DbDebugEnable bool `env:"DB_DEBUG_ENABLE" envDefault:"true"`
}

type Service struct {
	*service.BaseApp
	setting *extraSetting
}

func NewService() *Service {
	s := &Service{
		service.NewApp("Crawl data Service", "v1.0"),
		&extraSetting{},
	}
	_ = env.Parse(s.setting)

	// db := s.GetDB()

	// if s.setting.DbDebugEnable {
	// 	db = db.Debug()
	// }

	v1Api := s.Router.Group("/api/v1")
	v1Api.POST("/auth/sign-up", nil)

	// migration
	// migrate := handlers.NewMigrationHandler(db)
	// s.Router.POST("/internal/migrate", migrate.Migrate)
	return s
}
