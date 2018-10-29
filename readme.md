## SiteScanGo

### another subdomain brute tool written in Go

### Go编写的子域名枚举工具
---
```
Usage of C:\Users\30393\AppData\Local\Temp\___2588go_build_main_go.exe:
  -dict string
    	brute-dict location. (default "./dict/domain.txt")
  -domain string
    	domain to brute (default "baidu.com")
  -save string
    	where to save results (default "./baidu.com.csv")
  -third
    	get third-level subdomain (slow)
  -title
    	get website title (slow)
  -version
    	print program version
```


TODO
- [x] 标题获取
- [ ] ~~自动选择最佳DNS server(应该不需要)~~
- [ ] more search api
- [x] 结果导出csv
- [x] 多级子域名
