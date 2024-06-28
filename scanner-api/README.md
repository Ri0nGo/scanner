# Scanner-api
scanner-api 是一个可以检测主机地址是否存活，端口是否开发的接口程序，使用Gin作为web框架


## 1. 目录结构
```text
├─handler        # 处理请求 
├─lib            # 第三方包，这里使用的pro-bing拉取到本地，做了一点调整
│  └─pro-bing   
├─pkg            # 项目中存放的包
│  └─scanner
├─router         # 路由层
├─util           # 工具层
└─wiki           # 文档目录
    ├─docs
    └─img
```

## 2. 特性

- [x] 支持ip/域名 ping检测
- [x] 支持ip/域名 端口检测
- [x] 支持多主机，多网段ping和端口检测
- [x] 支持扫描结果定时清空
- [x] 支持中途停止扫描

# 3. 项目说明

1. 项目没有设计到数据库，所以没有创建svc，repo等目录，目录结构比较随意。
2. 代码中存在print函数，代码洁癖的大佬可以自己去掉。

# 4. WebUI 参考图

![WebUI.png](wiki%2Fimg%2FWebUI.png)

# 5. 扫描执行流程图

![WebUI.png](扫描执行流程图.png)
