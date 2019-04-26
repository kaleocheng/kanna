# Kanna

Kanna is a tiny tool inspired by nodemon. It watches the files in the directory in where Kanna started or by your specify, 
so if any files change, Kanna automatically restarts your command.


## Quick start guide

### Install
```bash
go get github.com/kaleocheng/kanna
```
### Usage
```bash
# simple run go build
kanna go build

# with flags, you can use ""
kanna "go build -v"

kanna node app.js
```


## Acknowledgements 

Special thanks go to [beego/bee](https://github.com/beego/bee), [nathany/looper](https://github.com/nathany/looper/blob/master/watch.go) and [gohugoio/hugo](https://github.com/gohugoio/hugo), from those projects, I gained a lot.
