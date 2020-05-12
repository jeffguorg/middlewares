package client

import (
	"github.com/go-redis/redis"
	"github.com/jeffguorg/middlewares/session"
)

type Client struct {
	rclient *redis.Client
}

func NewCluster(options *redis.Options) Client {
	redisClient := redis.NewClient(options)

	return Client{rclient: redisClient}
}

func (client Client) Load(name string) (map[string]interface{}, error) {
	panic("implement me")
}

func (client Client) Reset(name string) error {
	panic("implement me")
}

func (client Client) Update(name string, value map[string]interface{}) error {
	panic("implement me")
}

var (
	_ session.Client = Client{}
)
