// a core for sentry capture error

package logging

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func sentryLevel(lvl zapcore.Level) sentry.Level {
	switch lvl {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel:
		return sentry.LevelFatal
	case zapcore.PanicLevel:
		return sentry.LevelFatal
	case zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelFatal
	}
}

// SentryCoreConfig is a minimal set of parameters for Sentry Core.
type SentryCoreConfig struct {
	Tags              map[string]string
	DisableStacktrace bool
	Level             zapcore.Level
	FlushTimeout      time.Duration
	Hub               *sentry.Hub
}

// sentryCore the core for sentry
type sentryCore struct {
	client *sentry.Client
	cfg    *SentryCoreConfig
	zapcore.LevelEnabler
	flushTimeout time.Duration

	fields map[string]interface{}
}

// GetSentryClient return sentry client
func (c *sentryCore) GetSentryClient() *sentry.Client {
	return c.client
}

func (c *sentryCore) with(fs []zapcore.Field) *sentryCore {
	// Copy our map.
	m := make(map[string]interface{}, len(c.fields))
	for k, v := range c.fields {
		m[k] = v
	}

	// Add fields to an in-memory encoder.
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}

	// Merge the two maps.
	for k, v := range enc.Fields {
		m[k] = v
	}

	return &sentryCore{
		client:       c.client,
		cfg:          c.cfg,
		fields:       m,
		LevelEnabler: c.LevelEnabler,
	}
}

// With zap core interface
func (c *sentryCore) With(fs []zapcore.Field) zapcore.Core {
	return c.with(fs)
}

// Check zap core interface
func (c *sentryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.cfg.Level.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// Write zap core interface
func (c *sentryCore) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	clone := c.with(fs)

	event := sentry.NewEvent()
	event.Message = ent.Message
	event.Timestamp = ent.Time.Unix()
	event.Level = sentryLevel(ent.Level)
	event.Platform = "pink-lady"
	event.Extra = clone.fields
	event.Tags = c.cfg.Tags

	if !c.cfg.DisableStacktrace {
		trace := sentry.NewStacktrace()
		if trace != nil {
			event.Exception = []sentry.Exception{{
				Type:       ent.Message,
				Value:      ent.Caller.TrimmedPath(),
				Stacktrace: trace,
			}}
		}
	}

	hub := c.cfg.Hub
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	_ = c.client.CaptureEvent(event, nil, hub.Scope())

	// We may be crashing the program, so should flush any buffered events.
	if ent.Level > zapcore.ErrorLevel {
		c.client.Flush(c.flushTimeout)
	}
	return nil
}

// Sync zap core interface
func (c *sentryCore) Sync() error {
	c.client.Flush(c.flushTimeout)
	return nil
}

// NewSentryCore new a sentry core
func NewSentryCore(cfg SentryCoreConfig, sentryClient *sentry.Client) zapcore.Core {

	core := sentryCore{
		client:       sentryClient,
		cfg:          &cfg,
		LevelEnabler: cfg.Level,
		flushTimeout: 3 * time.Second,
		fields:       make(map[string]interface{}),
	}

	if cfg.FlushTimeout > 0 {
		core.flushTimeout = cfg.FlushTimeout
	}

	return &core
}

// SentryAttach attach sentrycore
func SentryAttach(l *zap.Logger, sentryClient *sentry.Client) *zap.Logger {
	cfg := SentryCoreConfig{
		Level: zap.ErrorLevel,
		Tags: map[string]string{
			"source": "zap sentry core",
		},
	}
	core := NewSentryCore(cfg, sentryClient)
	return AttachCore(l, core)
}
