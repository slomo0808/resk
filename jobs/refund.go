package jobs

import (
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/slomo0808/infra"
	"github.com/slomo0808/infra/comm"
	"imooc.com/resk/core/envelopes"
	"time"
)

type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
	mutex  *redsync.Mutex
}

func (s *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	d := ctx.Props().GetDurationDefault("jobs.refund.interval", 1*time.Minute)
	s.ticker = time.NewTicker(d)
	maxIdle := ctx.Props().GetIntDefault("redis.maxIdle", 2)
	maxActive := ctx.Props().GetIntDefault("redis.maxActive", 5)
	idleTimeout := ctx.Props().GetDurationDefault("redis.idleTimeout", 20*time.Second)
	addr := ctx.Props().GetDefault("redis.addr", "127.0.0.1:6379")
	pools := make([]redsync.Pool, 0)
	pool := &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", addr)
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
	pools = append(pools, pool)
	rsync := redsync.New(pools)
	ip := comm.GetIP()
	s.mutex = rsync.NewMutex("lock:RefundExpired",
		redsync.SetExpiry(50*time.Second),
		redsync.SetTries(3),
		redsync.SetGenValueFunc(func() (s string, e error) {
			now := time.Now()
			logrus.Infof("节点%s正在执行过期红包的退款业务", ip)
			return fmt.Sprintf("%d:%s", now.Unix(), ip), nil
		}))
}

func (s *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	go func() {
		for {
			c := <-s.ticker.C
			err := s.mutex.Lock()
			if err == nil {
				logrus.Debug("过期红包退款开始。。。", c)
				// 红包过期退款业务逻辑代码\
				domain := new(envelopes.ExpiredEnvelopeDomain)
				domain.Expired()
			} else {
				logrus.Info("已经有节点在运行该任务,err=", err.Error())
			}
			s.mutex.Unlock()
		}
	}()
}

func (s *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {
	s.ticker.Stop()
}
