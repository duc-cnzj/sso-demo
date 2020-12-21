package server

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"time"
)

type RedisLoader struct {
}

func (r *RedisLoader) GetWeight() int {
	return 3
}

func (r *RedisLoader) Load(s *Server) error {
	s.redis = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", s.config.RedisHost, s.config.RedisPort))

			if err != nil {
				log.Debug().Err(err).Msg("RedisLoader Load")
				return nil, err
			}

			if s.config.RedisPassword != "" {
				if _, err := c.Do("AUTH", s.config.RedisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
	}

	log.Info().Msg("Redis loaded.")

	return nil
}
