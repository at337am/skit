*构建:*  

windows:  
```bash
go build -ldflags="-H windowsgui" -o hello.exe main.go
```

linux:  
```bash
go build -o hello main.go
```

*clean:*  
```bash
/bin/rm -rfv middleware routes internal go.mod go.sum main.go README.md config/config.go
```

