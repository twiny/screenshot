package api

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/time/rate"
)

const (
	ChromeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36"
)

// Store
type Store interface {
	SaveImage(img *Image) error
	FindImage(key string) (*Image, error)
	Clean(d time.Duration) error
	Statistic() (map[string]interface{}, error)
	Close()
}

// Limiter
type Limiter interface {
	FindIPAddr(ip string) *rate.Limiter
	Clean(d time.Duration)
	Statistic() map[string]interface{}
}

// Logger
type Logger interface {
	Info(on, message string, properties map[string]string)
	Error(on, message string, properties map[string]string)
	Fetal(on, message string, properties map[string]string)
	Close()
}

// App
type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	//
	config  *Config
	store   Store
	limiter Limiter
	logger  Logger
	//
	shutdown chan struct{}
}

// NewApp
func NewApp(config *Config, store Store, limit Limiter, log Logger) *App {
	debug := !config.Debug

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserDataDir(config.ChromePath),
		chromedp.UserAgent(ChromeUserAgent),
		chromedp.Flag("headless", debug),
		chromedp.Flag("disable-gpu", debug),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &App{
		ctx:    ctx,
		cancel: cancel,
		//
		config:  config,
		store:   store,
		limiter: limit,
		logger:  log,
		//
		shutdown: make(chan struct{}),
	}
}

// Sync
func (a *App) Sync() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-a.shutdown:
				close(a.shutdown)
				return
			case <-ticker.C:
				// clean image
				if err := a.store.Clean(a.config.ImageCache); err != nil {
					a.logger.Error("Sync", err.Error(), map[string]string{
						"func": "Clean",
					})
				}

				// clean limiter
				// every 10 * min clean all limiters
				a.limiter.Clean(10 * time.Minute)
			}
		}
	}()
}

// Capture ScreenShot
func (a *App) Capture(link string, width, height int64) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(a.ctx)
	defer cancel()

	// run
	var screenshot []byte
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(width, height),
		chromedp.Navigate(link),
		chromedp.WaitReady("body"),
		chromedp.CaptureScreenshot(&screenshot),
	); err != nil {
		return nil, err
	}

	return screenshot, nil
}

// Close
func (a *App) Close() {
	a.shutdown <- struct{}{}
	a.store.Close()
	a.logger.Close()
}
