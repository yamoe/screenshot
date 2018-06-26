# Build

## Centos

1. install go : <https://golang.org/dl/>
    ```bash
    wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
    sudo tar -C /usr/local go1.10.3.linux-amd64.tar.gz
    vi ~/.bashrc
        export PATH=$PATH:/usr/local/go/bin
    ```

1. install dep : <https://golang.github.io/dep/>
    ```bash
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
        (기본으로 $GOPATH/bin 에 생성됨)
    sudo chown root.root dep
    sudo mv dep /usr/local/go/bin
    ```

1. init project
    ```bash
    export GOPATH=/home/ky/tmp/go
    cd $GOPATH
    ```

1. init dep
    ```bash
    mkdir -p src/test.com/example
    cd src/test.com/example
    dep init
        (Gopkg.lock, Gopkg.toml, vendor)
    ```

    * 의존성 등록
        ```bash
        dep ensure
        dep ensure -add github.com/labstack/echo
        ```
    * 상태 확인
        ```bash
        dep status
        ```

## Windows

1. install go : <https://golang.org/dl/>
    https://dl.google.com/go/go1.10.3.windows-amd64.msi

1. install dep : <https://github.com/golang/dep/releases>
    https://github.com/golang/dep/releases/download/v0.4.1/dep-windows-amd64.exe

    c:\Go\bin\dep.exe 로 복사

1. init project
    ```bash
    > setx GOPATH c:\project\sshot
    ```
    cmd 창 재시작

1. init dep
    ```bash
    > cd ../sshot/src/newera/sshot
    > dep ensure
    ```

## compile

```bash
cd $GOPATH/src/newera/sshot
go install
"""result : $GOPATH/bin/sshot"""
```


## Patch

EventLoadEventFired 이벤트를 기다리기 위하여 아래 Pull Request 소스를 수동 적용시킴

<https://github.com/chromedp/chromedp/pull/79>

golang dep 으로 의존성 받은 후 patch 폴더 소스 덮어쓰면 됨
cp -r ./patch/* ./src/newera/dashboard2png/vendor/

## vscode

    go extension 설치

    프로젝트 열고 오른쪽 아래 "Analysis Tools Missing" 클릭하여 설치

## reference

* https://developers.google.com/web/updates/2017/04/headless-chrome
* https://github.com/chromedp/chromedp
