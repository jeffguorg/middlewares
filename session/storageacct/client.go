package client

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/jeffguorg/middlewares/session"
)

// Client update or load sesson info into/from redis
type Client struct {
	rclient *redis.Client
	keyFmt  string
}

// New return a new Client instance
func New(keyFmt string, options *redis.Options) Client {
	redisClient := redis.NewClient(options)

	return Client{rclient: redisClient, keyFmt: keyFmt}
}

// Load load session info from redis
func (client Client) Load(name string) (map[string]interface{}, error) {
	result, err := client.rclient.HGetAll(fmt.Sprintf(client.keyFmt, name)).Result()
	if err != nil {
		return nil, err
	}
	session := make(map[string]interface{})
	for k, v := range result {
		session[k] = v
	}

	return session, nil
}

// Reset clears session info
func (client Client) Reset(name string) error {
	return client.rclient.Del(fmt.Sprintf(client.keyFmt, name)).Err()
}

// Update update info in session
func (client Client) Update(name string, value map[string]interface{}) error {
	return client.rclient.HMSet(fmt.Sprintf(client.keyFmt, name), value).Err()
}

var (
	_ session.Client = Client{}
)
