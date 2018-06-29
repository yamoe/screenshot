package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

// Size contains html size infomation
type Size struct {
	DevicePixelRatio int
	TotalWidth       int
	ViewportWidth    int
	TotalHeight      int
	ViewportHeight   int
}

// Init is set member variables
func (s *Size) Init(ctxt context.Context, cdp *chromedp.CDP) error {
	return cdp.Run(ctxt, chromedp.Tasks{
		chromedp.Evaluate("window.devicePixelRatio", &s.DevicePixelRatio),
		chromedp.Evaluate("document.body.parentNode.scrollWidth", &s.TotalWidth),
		chromedp.Evaluate("document.body.clientWidth", &s.ViewportWidth),
		chromedp.Evaluate("document.body.parentNode.scrollHeight", &s.TotalHeight),
		chromedp.Evaluate("window.innerHeight", &s.ViewportHeight),
	})
}

// ScreenshotParam is parameter for screenshooter
type ScreenshotParam struct {
	Timeout     int
	Debug       bool
	URL         string
	Width       int
	WaitSec     int
	WaitVisible []string
	Javascript  string
	Cookies     []*http.Cookie
	Filepath    string
	ChromePort int
}

// Screenshot capture web page
type Screenshot struct {
	cdp    *chromedp.CDP
	ctxt   context.Context
	cancel context.CancelFunc
}

// Init it init. need "defer Uninit()"
func (s *Screenshot) Init(timeout int, debug bool, port int) error {
	timeoutSec := time.Duration(timeout) * time.Second
	ctxt, cancel := context.WithTimeout(context.Background(), timeoutSec)

	logFunc := func(string, ...interface{}) {}
	if debug {
		logFunc = log.Printf
	}

	cdp, err := chromedp.New(
		ctxt, chromedp.WithRunnerOptions(
			runner.Flag("no-sandbox", true),
			runner.Flag("headless", true),
			runner.Flag("disable-gpu", true),
			runner.Flag("hide-scrollbars", true),
			runner.Flag("no-first-run", true),
			runner.Flag("no-default-browser-check", true),
			runner.Flag("remote-debugging-port", port),
		),
		chromedp.WithLog(logFunc),
	)

	if err != nil {
		cancel()
	} else {
		s.ctxt = ctxt
		s.cdp = cdp
		s.cancel = cancel
	}

	return err
}

// Uninit is uninit
func (s *Screenshot) Uninit() {
	if s.cdp != nil {
		s.cdp.Shutdown(s.ctxt)
		s.cancel()
		s.cdp.Wait()
		s.cdp = nil
	}
}

// Capture is that
func (s *Screenshot) Capture(p *ScreenshotParam) (buf *[]byte, err error) {

	if err = s.setCookie(p.URL, p.Cookies); err != nil {
		err = fmt.Errorf("set cookie- %v", err)
		return
	}

	if err = s.open(p.URL); err != nil {
		err = fmt.Errorf("open- %v", err)
		return
	}

	if err = s.resize(p.Width); err != nil {
		err = fmt.Errorf("resize- %v", err)
		return
	}

	if err = s.waitSec(p.WaitSec); err != nil {
		err = fmt.Errorf("waitsec- %v", err)
		return
	}

	for _, v := range p.WaitVisible {
		if err = s.waitVisible(v); err != nil {
			err = fmt.Errorf("waitvisible- %v", err)
			return
		}
	}

	if err = s.runjs(p.Javascript); err != nil {
		err = fmt.Errorf("runjs- %v", err)
		return
	}

	buf, err = s.capture()
	if err != nil {
		err = fmt.Errorf("capture- %v", err)
		return
	}
	return
}

func (s *Screenshot) run(action chromedp.Action) error {
	return s.cdp.Run(s.ctxt, action)
}

func (s *Screenshot) setCookie(urlstr string, cookies []*http.Cookie) error {
	u, err := url.Parse(urlstr)
	if err != nil {
		return err
	}
	hostname := u.Hostname()

	return s.run(chromedp.Tasks{
		chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
			for _, c := range cookies {
				success, err := network.SetCookie(c.Name, c.Value).
					WithDomain(hostname).
					Do(ctxt, h)

				if err != nil {
					return err
				}
				if !success {
					return fmt.Errorf("could not set cookie : %s,%s", c.Name, c.Value)
				}
			}
			return nil
		}),
	})
}

func (s *Screenshot) open(urlstr string) error {
	return s.run(chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitEventLoad(),
	})
}

func (s *Screenshot) waitSec(sec int) error {
	second := time.Duration(sec) * time.Second
	return s.run(chromedp.Tasks{
		chromedp.Sleep(second),
	})
}

func (s *Screenshot) waitVisible(sel string) error {
	if sel == "" {
		return nil
	}
	return s.run(chromedp.Tasks{
		chromedp.WaitVisible(sel, chromedp.ByQueryAll),
	})
}

func (s *Screenshot) resize(width int) error {
	size := new(Size)
	if err := size.Init(s.ctxt, s.cdp); err != nil {
		return err
	}

	var w int64

	if width > 0 {
		w = int64(width)
	} else {
		w = int64(size.TotalWidth)
	}

	h := int64(size.ViewportHeight)

	if err := s.cdp.Run(s.ctxt, emulation.SetDeviceMetricsOverride(w, h, 0, false)); err != nil {
		return err
	}

	// repeat one more
	if err := size.Init(s.ctxt, s.cdp); err != nil {
		return err
	}
	w = int64(size.TotalWidth)
	h = int64(size.TotalHeight)

	return s.cdp.Run(s.ctxt, emulation.SetDeviceMetricsOverride(w, h, 0, false))
}

func (s *Screenshot) capture() (*[]byte, error) {
	var buf []byte
	err := s.run(chromedp.CaptureScreenshot(&buf))
	return &buf, err
}

func (s *Screenshot) runjs(js string) error {
	js = strings.TrimSpace(js)
	if len(js) <= 0 {
		return nil
	}

	var buf []byte
	return s.run(chromedp.Evaluate(js, &buf))
}

// RunScreenshot is web page sccrenshot
func RunScreenshot(p *ScreenshotParam) error {

	ss := new(Screenshot)
	if err := ss.Init(p.Timeout, p.Debug, p.ChromePort); err != nil {
		return err
	}
	defer ss.Uninit()

	buf, err := ss.Capture(p)
	if err != nil {
		return err
	}
	return SaveFile(buf, p.Filepath)
}

// SaveFile is save file
func SaveFile(buf *[]byte, path string) error {
	return ioutil.WriteFile(path, *buf, 0644)
}
