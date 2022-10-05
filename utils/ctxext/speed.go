// Package ctxext ctx扩展
package ctxext

import (
	"time"
	"unsafe"

	rei "github.com/fumiama/ReiBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

// DefaultSingle 默认反并发处理
//
//	按 qq 号反并发
//	并发时返回 "您有操作正在执行, 请稍后再试!"
var DefaultSingle = rei.NewSingle(
	rei.WithKeyFn(func(ctx *rei.Ctx) int64 {
		return ctx.Message.From.ID
	}),
	rei.WithPostFn[int64](func(ctx *rei.Ctx) {
		_, _ = ctx.SendPlainMessage(false, "您有操作正在执行, 请稍后再试!")
	}),
)

// defaultLimiterManager 默认限速器管理
//
//	每 10s 5次触发
var defaultLimiterManager = rate.NewManager[int64](time.Second*10, 5)

type fakeLM struct {
	limiters unsafe.Pointer
	interval time.Duration
	burst    int
}

// SetDefaultLimiterManagerParam 设置默认限速器参数
//
//	每 interval 时间 burst 次触发
func SetDefaultLimiterManagerParam(interval time.Duration, burst int) {
	f := (*fakeLM)(unsafe.Pointer(defaultLimiterManager))
	f.interval = interval
	f.burst = burst
}

// LimitByUser 默认限速器 每 10s 5次触发
//
//	按 qq 号限制
func LimitByUser(ctx *rei.Ctx) *rate.Limiter {
	return defaultLimiterManager.Load(ctx.Message.From.ID)
}

// LimitByGroup 默认限速器 每 10s 5次触发
//
//	按群号限制
func LimitByGroup(ctx *rei.Ctx) *rate.Limiter {
	return defaultLimiterManager.Load(ctx.Message.Chat.ID)
}

// LimiterManager 自定义限速器管理
type LimiterManager struct {
	m *rate.LimiterManager[int64]
}

// NewLimiterManager 新限速器管理
func NewLimiterManager(interval time.Duration, burst int) (m LimiterManager) {
	m.m = rate.NewManager[int64](interval, burst)
	return
}

// LimitByUser 自定义限速器
//
//	按 qq 号限制
func (m LimiterManager) LimitByUser(ctx *rei.Ctx) *rate.Limiter {
	return m.m.Load(ctx.Message.From.ID)
}

// LimitByGroup 自定义限速器
//
//	按群号限制
func (m LimiterManager) LimitByGroup(ctx *rei.Ctx) *rate.Limiter {
	return m.m.Load(ctx.Message.Chat.ID)
}
