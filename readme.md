## SiteScanGo

### not only a subdomain brute tool written in Go 

### Go编写的不只是子域名枚举的工具
---
```
  -dict string
    	brute-dict location. default ./dict/domain.txt (default "./dict/domain.txt")
  -domain string
    	determine target  (default "baidu.com")
  -sendir
    	brute sensitive directory of target
  -sub
    	brute subdomains of target
  -third
    	get third-level info (slow)
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
- [x] 敏感目录扫描
- [ ] web界面(probably not)