package api

import (
	"encoding/json"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
)

// limiter
func (mux *Server) limiter(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// get client ip addr
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			mux.app.logger.Info("limiter", err.Error(), map[string]string{
				"func": "SplitHostPort",
			})
		}
		// get rate
		rate := mux.app.limiter.FindIPAddr(ip)

		// check rate
		if !rate.Allow() {
			mux.writer(w, http.StatusTooManyRequests, "slow down")
			return
		}

		// next
		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)
}

// takeScreenShot
func (mux *Server) takeScreenShot(w http.ResponseWriter, r *http.Request) {
	data := struct {
		URL    string `json:"url"`
		Width  int64  `json:"width"`
		Height int64  `json:"height"`
	}{}

	// json decoder
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&data); err != nil {
		mux.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	if data.Width == 0 {
		data.Width = 1024
	}

	if data.Height == 0 {
		data.Height = 800
	}

	// img
	img := &Image{
		UUID:      xid.New().String(),
		Status:    ImageStatusPending,
		Message:   "",
		Binary:    nil,
		CreatedAt: time.Now(),
	}

	// save img
	if err := mux.app.store.SaveImage(img); err != nil {
		mux.app.logger.Error("takeScreenShot", err.Error(), map[string]string{
			"func": "SaveImage",
		})
	}

	go func(i *Image) {
		data, err := mux.app.Capture(data.URL, data.Width, data.Height)
		if err != nil {
			// save image error
			img.Status = ImageStatusFail
			img.Message = err.Error()

			if err := mux.app.store.SaveImage(img); err != nil {
				mux.app.logger.Error("goroutine", err.Error(), map[string]string{
					"func": "SaveImage",
				})
			}

			mux.app.logger.Error("goroutine", err.Error(), map[string]string{
				"func": "Capture",
			})
			return
		}

		// save img success
		img.Status = ImageStatusSuccess
		img.Message = ""
		img.Binary = data

		if err := mux.app.store.SaveImage(img); err != nil {
			mux.app.logger.Error("goroutine", err.Error(), map[string]string{
				"func": "SaveImage",
			})
		}
	}(img)

	mux.writer(w, http.StatusOK, struct {
		UUID string `json:"uuid"`
	}{
		UUID: img.UUID,
	})
}

// findScreenShot
func (mux *Server) findScreenShot(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		mux.writer(w, http.StatusInternalServerError, "invalid uuid")
		return
	}

	// find screenshot
	img, err := mux.app.store.FindImage(uuid)
	if err != nil {
		mux.writer(w, http.StatusNotFound, err.Error())
		return
	}

	// if image fail
	if img.Status == ImageStatusFail {
		mux.writer(w, http.StatusInternalServerError, img.Message)
		return
	}

	// if image pending
	if img.Status == ImageStatusPending {
		mux.writer(w, http.StatusInternalServerError, "screenshot not yet captured.")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(img.Binary)))

	w.Write(img.Binary)
}

// stats
func (mux *Server) stats(w http.ResponseWriter, r *http.Request) {
	sstats, err := mux.app.store.Statistic()
	if err != nil {
		mux.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	memstats := new(runtime.MemStats)
	runtime.ReadMemStats(memstats)

	// astats
	astats := map[string]interface{}{
		"Goroutine":   runtime.NumGoroutine(),
		"CPUs":        runtime.NumCPU(),
		"System":      PrintBytes(memstats.Sys),
		"StackSystem": PrintBytes(memstats.StackSys),
		"GCSize":      PrintBytes(memstats.GCSys),
		"CompletedGC": memstats.NumGC,
	}

	data := map[string]interface{}{
		"Store":     sstats,
		"Limiter":   mux.app.limiter.Statistic(),
		"Resources": astats,
	}

	mux.writer(w, http.StatusOK, data)
}
