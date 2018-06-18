# Screenshot Java 예제

## Prerequisite

### Common

* java 8 (64bit)

### Windows

* chrome : https://www.google.com/chrome/
* chromdriver : https://sites.google.com/a/chromium.org/chromedriver/downloads

    chromdriver는 다운로드 위 OS PATH에 등록

### CentOS 7

* sudo yum -y install epel-release

* sudo yum -y install java chromedriver chromium

* sudo yum -y install baekmuk-ttf-dotum-fonts

    http://naver.com/ 캡쳐시 한글 깨지므로 돋움 font 

* 테스트한 버전

    ChromeDriver 2.35 (0)

    Chromium 65.0.3325.181 Fedora Project
    
    (xorg-x11-server-Xvfb 없이도 정상동작하여 pass)

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



