# 说明

 	此工具处理漏洞扫描ip资产端口开放过多的数据

​	在渗透测试中会出现一个ip因为蜜罐、安全组端口全开等问题，导致很多资产都会出现这种情况，在后续的自动化扫描会导致运行超时、扫描无效应用等问题，这里采取数据清洗进行优化。

# 使用

```
ipHandle -input {{Output}}/portscan/open-ports.txt -output {{Output}}/portscan/open-ports-handle.txt
```

​	默认会删除一个ip开放100个端口以上的资产

**本工具会继承到本人自己开发的全自动化工具流漏洞挖掘系统中**
