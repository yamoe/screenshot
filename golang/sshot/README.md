# SSHOT

## Prerequisite

chrome 설치

* Centos
    ```bash
    sudo yum -y install epel-release
    sudo yum -y install chromium

    # web page font
    sudo yum -y install baekmuk-ttf-dotum-fonts
    ```

* windows
   <https://www.google.com/chrome/>

## Usage

* help

    ```text
    $ ./sshot -h
    Options:

    -h, --help            display help information
    -r, --url            *web url
    -f, --filepath[=sshot.png]   png filepath for save
    -w, --width[=0]       web width
    -s, --waitsec[=0]     wait seconds(sleep)
    -v, --waitvisible     wait visible css
    -j, --javascript      execute javascript
    -t, --timeout[=120]   timeout
    -l, --loginurl        login url
    -u, --username        username
    -p, --password        password
    -d, --debug[=false]   print debug log
    ```

* screenshot - naver.com 

    ./sshot -r "http://naver.com" -f naver.png
