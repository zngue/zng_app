package db

import (
	"context"
	"echo_layout_v2/internal/conf"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
)

type RedisStream struct{}

type StreamMsg struct {
	Id  string
	Val *MsgVal
}

// stream会强制转换类型string的json，所以每个属性必须为string
type MsgVal struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
	Time string `json:"time"`
}

func NewRedisStream() *RedisStream {
	return &RedisStream{}
}

func (rs *RedisStream) structToMap(msg *MsgVal) (map[string]interface{}, error) {
	tmpMap := make(map[string]interface{})
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &tmpMap)
	if err != nil {
		return nil, err
	}
	return tmpMap, nil
}

func (rs *RedisStream) mapToStruct(tmpMap map[string]interface{}) (*MsgVal, error) {
	msg := &MsgVal{}
	msgBytes, err := json.Marshal(tmpMap)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// CreateStreamAndGroup 第一次创建stream和group（只允许第一次）
func (rs *RedisStream) CreateStreamAndGroup(ctx context.Context) error {
	config, err := conf.GetConfig()
	if err != nil {
		return err
	}
	count, err := rdb.Exists(ctx, config.Stream.TopicKey).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		// 没有被创建
		tmpMap, err := rs.structToMap(&MsgVal{
			Name: "first",
			Pwd:  "new",
		})
		if err != nil {
			return err
		}

		id, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: config.Stream.TopicKey,
			MaxLen: config.Stream.MaxLen,
			Values: tmpMap,
		}).Result()
		if err != nil {
			return err
		}
		// 删除id的数据
		_, err = rdb.XDel(ctx, config.Stream.TopicKey, id).Result()
		if err != nil {
			return err
		}
	}

	// 创建group
	_, err = rdb.XGroupCreate(ctx, config.Stream.TopicKey, config.Stream.GroupName, "$").Result()
	if err != nil {
		return err
	}

	return nil
}

// CreateMsg 向对应stream写消息
func (rs *RedisStream) CreateMsg(ctx context.Context, msg *MsgVal) (string, error) {
	config, err := conf.GetConfig()
	if err != nil {
		return "", err
	}

	msgMap, err := rs.structToMap(msg)
	if err != nil {
		return "", err
	}

	id, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: config.Stream.TopicKey,
		MaxLen: config.Stream.MaxLen,
		Values: msgMap,
	}).Result()
	if err != nil {
		return "", err
	}
	return id, nil
}

// ReadGroupMsg 读取一条stream的消息id和内容
func (rs *RedisStream) ReadGroupMsg(ctx context.Context) (*StreamMsg, error) {
	config, err := conf.GetConfig()
	if err != nil {
		return nil, err
	}
	st, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    config.Stream.GroupName,
		Consumer: config.Stream.ConsumerName,
		Block:    -1,
		Streams:  []string{config.Stream.TopicKey, ">"},
		Count:    1,
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	if len(st) <= 0 {
		return nil, nil
	}

	msg, err := rs.mapToStruct(st[0].Messages[0].Values)
	if err != nil {
		return nil, err
	}

	sm := &StreamMsg{
		Id:  st[0].Messages[0].ID,
		Val: msg,
	}
	return sm, nil
}

// ReadPendingMsg 读取一条等待中未发送ACK的消息
func (rs *RedisStream) ReadPendingMsg(ctx context.Context) (*StreamMsg, error) {
	config, err := conf.GetConfig()
	if err != nil {
		return nil, err
	}

	// 先获取pending的消息ID
	pe, err := rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: config.Stream.TopicKey,
		Group:  config.Stream.GroupName,
		Start:  "-",
		End:    "+",
		Count:  1,
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	if len(pe) <= 0 {
		return nil, nil
	}

	// 获取消息体数据
	peMsg, err := rdb.XRange(ctx, config.Stream.TopicKey, pe[0].ID, pe[0].ID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	if len(peMsg) <= 0 {
		return nil, nil
	}

	msg, err := rs.mapToStruct(peMsg[0].Values)
	if err != nil {
		return nil, err
	}

	sm := &StreamMsg{
		Id:  peMsg[0].ID,
		Val: msg,
	}
	return sm, err
}

func (rs *RedisStream) SetACK(ctx context.Context, id string) error {
	config, err := conf.GetConfig()
	if err != nil {
		return err
	}
	_, err = rdb.XAck(ctx, config.Stream.TopicKey, config.Stream.GroupName, id).Result()
	if err != nil {
		return err
	}
	return nil
}