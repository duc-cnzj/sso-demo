package server

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"time"
)

type DBLoader struct {
}

func (c *DBLoader) GetWeight() int {
	return 2
}

func (c *DBLoader) Load(s *Server) error {
	var err error
	s.db, err = gorm.Open(s.config.DBConnection, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci", s.config.DBUsername, s.config.DBPassword, s.config.DBHost, s.config.DBPort, s.config.DBDatabase))
	if err != nil {
		log.Debug().Err(err).Msg("DBLoader Load")
		return err
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	s.db.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	s.db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	s.db.DB().SetConnMaxLifetime(time.Hour)

	log.Info().Msg("DB loaded.")

	return nil
}
