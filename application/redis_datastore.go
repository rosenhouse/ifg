package application

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type RedisDataStore struct {
	Host     string
	Password string
	pool     *redis.Pool
}

func (s *RedisDataStore) checkConnection() error {
	c, err := redis.Dial("tcp", s.Host)
	if err != nil {
		return err
	}
	_, err = c.Do("AUTH", s.Password)
	if err != nil {
		c.Close()
		return err
	}
	c.Close()
	return nil
}

func (s *RedisDataStore) Initialize() error {
	err := s.checkConnection()
	if err != nil {
		return err
	}
	s.pool = &redis.Pool{
		MaxIdle:   5,
		MaxActive: 5,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", s.Host)
			if err != nil {
				return nil, err
			}
			_, err = c.Do("AUTH", s.Password)
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	return nil
}

func (s *RedisDataStore) Get(key string) ([]byte, error) {
	c := s.pool.Get()
	defer c.Close()

	val, err := c.Do("GET", key)
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
	c := s.pool.Get()
	defer c.Close()

	_, err := c.Do("SET", key, val)
	return err
}
