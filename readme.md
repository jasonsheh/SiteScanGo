## SiteScanGo

### not only a subdomain brute tool written in Go 

### Go编写的不只是子域名枚举的工具
---
```
  Usage of SiteScanGo:
    -cclass
      	brute c class of ips
    -dict string
      	brute-dict location. default ./dict/domain.txt (default "./dict/domain.txt")
    -dir
      	brute sensitive directory of target
    -port
      	get ip open port only work with sub domain brute otherwIse use nmap or masscan
    -sub
      	brute subdomains of target
    -target string
      	determine target  (default "baidu.com")
    -third
      	get third-level info (slow)
    -title
      	get website title (slow)
    -version
      	print program version

```

---
dns服务器的选取会影响最终结果

本地测试 119.29.29.29 好于 223.6.6.6
---

TODO
- [x] 标题获取
- [ ] more search api
- [x] 结果导出csv
- [x] 多级子域名
- [x] 敏感目录扫描
- [ ] C段扫描