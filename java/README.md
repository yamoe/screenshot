# Screenshot Java 예제

## Prerequisite

### Common

* java 8 (64bit)

### Windows

* chrome : https://www.google.com/chrome/
* chromdriver : https://sites.google.com/a/chromium.org/chromedriver/downloads

    chromdriver는 다운로드 위 OS PATH에 등록

### CentOS 7

* yum -y install epel-release

* yum -y install chromedriver chromium xorg-x11-server-Xvfb

    ChromeDriver 2.35 (0)

    Chromium 65.0.3325.181 Fedora Project

    xorg-x11-server-Xvfb-1.19.5-5.el7.x86_64

* 캡쳐화면의 한글이 깨지는 경우 폰트 설치 필요

    yum -y install fonts-nanum-coding


## Run

### web page

1. java -jar Screenshot-0.0.1-SNAPSHOT.jar
1. http://localhost:8080

### call url

* http://localhost:8080/screenshot?url=http://naver.com

    캡쳐화면 보기

* http://localhost:8080/screenshot?url=http://naver.com&width=2000

   width 2000 px 로 캡쳐 

* http://localhost:8080/screenshot?url=http://naver.com&filename=image.png

    캡쳐화면 image.png 파일로 다운로드

* http://localhost:8080/screenshot?url=http://naver.com&tag=true

    캡쳐 화면을 html에 embed 할 수 있는 <img> 태그 출력

    <img src="data:image/png;base64,iVBORw0KGgoAAAANSUh.......>



