package redis

import (
	"time"

	"github.com/guestin/kboot"
	"github.com/guestin/log"
	"github.com/pkg/errors"
	goRedis "github.com/redis/go-redis/v9"
)

//goland:noinspection ALL
var logger log.ClassicLog

//goland:noinspection ALL
var zapLogger log.ZapLog

func init() {
	kboot.RegisterUnit("redis", _init)
}

func _init(unit kboot.Unit) (kboot.ExecFunc, error) {
	logger = kboot.GetTaggedLogger(ModuleName)
	zapLogger = kboot.GetTaggedZapLogger(ModuleName)
	_cfg = new(Config)
	err := kboot.UnmarshalSubConfig(ModuleName, _cfg,
		kboot.MustBindEnv(cfgKeyHost),
		kboot.MustBindEnv(cfgKeyPort),
		kboot.MustBindEnv(cfgKeyDbIdx),
		kboot.MustBindEnv(cfgKeyPassword),
		kboot.MustBindEnv(cfgKeyKeyPrefix),
	)
	if err != nil {
		return nil, err
	}
	if _cfg.KeyPrefix == "" {
		_cfg.KeyPrefix = DefaultKeyPrefix
	}
	option := &goRedis.UniversalOptions{
		Addrs:           _cfg.Address,
		DB:              _cfg.DbIdx,
		Password:        _cfg.Password,
		PoolSize:        _cfg.PoolSize,
		PoolTimeout:     time.Duration(_cfg.PoolTimeoutSec) * time.Second,
		MinIdleConns:    _cfg.MinIdleConns,
		MaxIdleConns:    _cfg.MaxIdleConns,
		MaxActiveConns:  _cfg.MaxActiveConns,
		ConnMaxIdleTime: time.Duration(_cfg.ConnMaxIdleTimeSec) * time.Second,
		ConnMaxLifetime: time.Duration(_cfg.ConnMaxLifetimeSec) * time.Second,
	}
	cli := goRedis.NewUniversalClient(option)
	if err = cli.Ping(unit.GetContext()).Err(); err != nil {
		return nil, errors.Wrap(err, "redis connect failed")
	}
	_redisCli = &Client{UniversalClient: cli}
	return _execute, nil
}

func _execute(unit kboot.Unit) kboot.ExitResult {
	<-unit.Done()
	return kboot.ExitResult{
		Code:  0,
		Error: nil,
	}
}
