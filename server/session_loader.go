package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/rs/zerolog/log"
)

type SessionLoader struct {
}

func (sl *SessionLoader) GetWeight() int {
	return 4
}

func (sl *SessionLoader) Load(s *Server) error {
	var (
		store redis.Store
		err   error
	)
	if store, err = redis.NewStoreWithPool(s.redis, []byte("secret")); err != nil {
		log.Error().Err(err).Msg("redis.NewStoreWithPool")
		return err
	}

	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: s.config.SessionLifetime,
	})
	s.session = store

	return nil
}
