package application

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type RedisDataStore struct {
	Host     string
	Password string
	conn     redis.Conn
}

func (s *RedisDataStore) Connect() error {
	c, err := redis.Dial("tcp", s.Host)
	if err != nil {
		return err
	}
	_, err = c.Do("AUTH", s.Password)
	if err != nil {
		c.Close()
		return err
	}
	s.conn = c
	return nil
}

func (s *RedisDataStore) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *RedisDataStore) Get(key string) ([]byte, error) {
	val, err := s.conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	bytes, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("unable to type-assert response of type %T to []byte", val)
	}
	return bytes, nil
}

func (s *RedisDataStore) Set(key string, val []byte) error {
	_, err := s.conn.Do("SET", key, val)
	return err
}
