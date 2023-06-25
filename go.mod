module github.com/saymoon/chatgpt-survivor

go 1.19

require (
	github.com/atotto/clipboard v0.1.4
	github.com/cherish-chat/chatgpt-firefox v0.0.3
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/playwright-community/playwright-go v0.2000.1
	github.com/sirupsen/logrus v1.9.3
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

replace github.com/cherish-chat/chatgpt-firefox => ./
