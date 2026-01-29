package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	capture "github.com/chattcp/chattcp-capture"
	"github.com/gin-gonic/gin"
)

type captureSession struct {
	cancel context.CancelFunc
	done   chan struct{}
}

var (
	captureMu      sync.Mutex
	currentSession *captureSession
)

func ListInterfaces(c *gin.Context) {
	interfaces, err := capture.ListNetworkInterfaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": interfaces,
	})
}

func StartCaptureSSE(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	done := make(chan struct{})
	session := &captureSession{cancel: cancel, done: done}
	captureMu.Lock()
	if currentSession != nil {
		currentSession.cancel()
		captureMu.Unlock()
		<-currentSession.done
		captureMu.Lock()
	}
	currentSession = session
	captureMu.Unlock()
	defer func() {
		captureMu.Lock()
		if currentSession == session {
			currentSession = nil
		}
		close(done)
		captureMu.Unlock()
	}()
	filter, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if filter.InterfaceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "interface required",
		})
		return
	}
	// SSE header
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	// SSE Listener
	listener := &SSEListener{
		ctx:    ctx,
		writer: c.Writer,
		done:   make(chan struct{}),
	}
	// Start capture
	captureErr := make(chan error, 1)
	go func() {
		defer close(listener.done)
		if err = capture.StartCapture(filter, listener); err != nil {
			select {
			case captureErr <- err:
			case <-ctx.Done():
				// 被新请求取消，忽略错误
			}
		}
	}()
	select {
	case <-ctx.Done():
		capture.StopCapture()
		<-listener.done
	case err = <-captureErr:
		listener.sendError(err.Error())
		<-listener.done
	}
}

type SSEListener struct {
	ctx    context.Context
	writer http.ResponseWriter
	done   chan struct{}
}

func (l *SSEListener) OnPkg(pkg *capture.OutputPacket) {
	select {
	case <-l.ctx.Done():
		return
	default:
	}
	data, err := json.Marshal(pkg)
	if err != nil {
		return
	}
	l.sendSSE("packet", string(data))
}

func (l *SSEListener) OnClose() {
	l.sendSSE("close", "stop capture")
}

func (l *SSEListener) sendSSE(event, data string) {
	select {
	case <-l.ctx.Done():
		return
	default:
		// event: <event>\ndata: <data>\n\n
		_, err := l.writer.Write([]byte("event: " + event + "\n"))
		if err != nil {
			return
		}
		_, err = l.writer.Write([]byte("data: " + data + "\n\n"))
		if err != nil {
			return
		}
		if flusher, ok := l.writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

func (l *SSEListener) sendError(errMsg string) {
	l.sendSSE("error", errMsg)
}

func parseFilterParams(c *gin.Context) (capture.FilterParam, error) {
	tcp := c.Query("tcp") == "true"
	udp := c.Query("udp") == "true"
	proto := ""
	if tcp {
		proto = "tcp"
	} else if udp {
		proto = "udp"
	}
	filter := capture.FilterParam{
		InterfaceName: c.Query("i"),
		Proto:         proto,
	}
	if hostSrc := c.Query("h.src"); hostSrc != "" {
		filter.Host = &struct {
			Src string
			Dst string
			Any string
		}{
			Src: hostSrc,
		}
	} else if hostDst := c.Query("h.dst"); hostDst != "" {
		filter.Host = &struct {
			Src string
			Dst string
			Any string
		}{
			Dst: hostDst,
		}
	} else if hostAny := c.Query("h"); hostAny != "" {
		filter.Host = &struct {
			Src string
			Dst string
			Any string
		}{
			Any: hostAny,
		}
	}
	if portSrcStr := c.Query("p.src"); portSrcStr != "" {
		portSrc, err := strconv.Atoi(portSrcStr)
		if err != nil {
			return filter, err
		}
		filter.Port = &struct {
			Src int
			Dst int
			Any int
		}{
			Src: portSrc,
		}
	} else if portDstStr := c.Query("p.dst"); portDstStr != "" {
		portDst, err := strconv.Atoi(portDstStr)
		if err != nil {
			return filter, err
		}
		filter.Port = &struct {
			Src int
			Dst int
			Any int
		}{
			Dst: portDst,
		}
	} else if portAnyStr := c.Query("p"); portAnyStr != "" {
		portAny, err := strconv.Atoi(portAnyStr)
		if err != nil {
			return filter, err
		}
		filter.Port = &struct {
			Src int
			Dst int
			Any int
		}{
			Any: portAny,
		}
	}
	return filter, nil
}
