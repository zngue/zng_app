package db

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client

type Stream struct {
	TopicKey     string
	MaxLen       int64
	GroupName    string
	Group        string
	Consumer     string
	Streams      []string // list of streams and ids, e.g. stream1 stream2 id1 id2
	Count        int64
	Block        time.Duration
	NoAck        bool
	ConsumerName string
}
type Config struct {
	Stream *Stream
}
type StreamMsg struct {
	Id   string
	Data string
}

// MsgVal stream会强制转换类型string的json，所以每个属性必须为string
type MsgVal struct {
	Data string `json:"data"`
}
type RedisStream struct {
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
func (rs *RedisStream) CreateStreamAndGroup(ctx context.Context) (err error) {
	config := &Config{
		Stream: &Stream{
			TopicKey: "test_stream",
		},
	}
	count, err := rdb.Exists(ctx, config.Stream.TopicKey).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		// 没有被创建
		tmpMap, err := rs.structToMap(&MsgVal{
			Data: "first",
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
func (rs *RedisStream) CreateMsg(ctx context.Context, data string) (string, error) {
	config := &Config{
		Stream: &Stream{
			TopicKey: "test_stream",
		},
	}
	// 把data打包成base64，避免怪异的符号影响json
	b64 := base64.StdEncoding.EncodeToString([]byte(data))
	msgMap, err := rs.structToMap(&MsgVal{
		Data: b64,
	})
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
	config := &Config{
		Stream: &Stream{
			TopicKey: "test_stream",
		},
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
	// 从base64解开
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return nil, err
	}
	sm := &StreamMsg{
		Id:   st[0].Messages[0].ID,
		Data: string(data),
	}
	return sm, nil
}

// ReadPendingMsg 读取一条等待中未发送ACK的消息
func (rs *RedisStream) ReadPendingMsg(ctx context.Context) (*StreamMsg, error) {
	config := &Config{
		Stream: &Stream{
			TopicKey:     "test_stream",
			ConsumerName: "consumer1",
			Block:        time.Second * 5,
			NoAck:        true,
			Count:        1,
			Group:        "group1",
		},
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
	// 从base64解开
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return nil, err
	}
	sm := &StreamMsg{
		Id:   peMsg[0].ID,
		Data: string(data),
	}
	return sm, err
}

// SetACK 设置消息已消费
func (rs *RedisStream) SetACK(ctx context.Context, id string) (err error) {
	config := &Config{
		Stream: &Stream{
			TopicKey: "test_stream",
		},
	}
	_, err = rdb.XAck(ctx, config.Stream.TopicKey, config.Stream.GroupName, id).Result()
	if err != nil {
		return err
	}
	return nil
}
