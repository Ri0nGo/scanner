# scanner
端口扫描器

基于协程池实现的端口扫描

## 特性
- [x] 自定义worker数量
- [x] 自定义端口
- [x] tcp 端口扫描


## 参数说明
```
D:\projects\scanner>scanner.exe --help
Usage of main.exe:
  -h string  设置主机地址
        host address (default "127.0.0.1")
  -p string  设置端口
        port range, use ps mode (default "1-1024")
  -type string  设置扫描类型，暂时仅支持tcp port 扫描
        scanner mode, tcp port scan(ps) or url address scan(us) (default "portScan")
  -url string   设置url地址
        url address use us mode (default "http://127.0.0.1")
  -w int    设置worker 数量，最多1000
        goroutine number, max 1000 (default 1)
```

## 使用
```
D:\projects\scanner>scanner.exe -p 80,445,102,3389,8000,8000 -w 2  
port: 445, status: open
port: 102, status: open
port: 80, status: close
port: 3389, status: close
port: 8000, status: close
port: 8000, status: close
use time: 4.0504124s
```

## 缺陷

- **待完善url地址扫描**
- **系统设计重构，使支持多主机端口扫描**
