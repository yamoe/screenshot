package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mkideal/cli"
)

const Version = "0.0.1"

type argT struct {
	cli.Helper
	Version     bool   `cli:"!v,version" usage:"print version"`
	URL         string `cli:"*r,url" usage:"web url"`
	Filepath    string `cli:"f,filepath" usage:"png filepath for save" dft:"sshot.png"`
	Width       int    `cli:"w,width" usage:"web width" dft:"0"`
	WaitSec     int    `cli:"s,waitsec" usage:"wait seconds(sleep)" dft:"0"`
	WaitVisible string `cli:"i,waitvisible" usage:"wait visible css"`
	Javascript  string `cli:"j,javascript" usage:"execute javascript"`
	Timeout     int    `cli:"t,timeout" usage:"timeout" dft:"120"`
	LoginURL    string `cli:"l,loginurl" usage:"login url"`
	Username    string `cli:"u,username" usage:"username"`
	Password    string `cli:"p,password" usage:"password"`
	Debug       bool   `cli:"d,debug" usage:"print debug log" dft:"false"`
	ChromePort  int    `cli:"o,port" usage:"chrome remote debugging port" dft:"9222"`
}

func needLogin(argv *argT) bool {
	return argv.Username != "" && argv.Password != "" && argv.LoginURL != ""
}

func login(urlstr string, username string, password string) []*http.Cookie {
	csrf, cookies, err := LoginData(urlstr)
	if err != nil {
		log.Panic(err)
	}

	cookies, err = Login(urlstr, username, password, csrf, cookies)
	if err != nil {
		log.Panic(err)
	}
	return cookies
}

func main() {
	var argv *argT
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv = ctx.Argv().(*argT)
		if argv.Version {
			ctx.String(fmt.Sprintf("%s\n", Version))
		}
		return nil
	})

	if argv == nil || argv.URL == "" {
		return
	}

	// login
	var cookies []*http.Cookie
	useLogin := needLogin(argv)
	if useLogin {
		cookies = login(argv.LoginURL, argv.Username, argv.Password)
	}

	// screenshot
	param := ScreenshotParam{
		Timeout:    argv.Timeout,
		Debug:      argv.Debug,
		URL:        argv.URL,
		Width:      argv.Width,
		WaitSec:    argv.WaitSec,
		Javascript: argv.Javascript,
		Cookies:    cookies,
		Filepath:   argv.Filepath,
		ChromePort: argv.ChromePort,
	}
	if argv.WaitVisible != "" {
		param.WaitVisible = append(param.WaitVisible, argv.WaitVisible)
	}

	if err := RunScreenshot(&param); err != nil {
		log.Panic(err)
	}

	os.Exit(0)
}
