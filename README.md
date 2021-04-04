# Log

基础于uber的高性能日志zap。已封装对控制台以及文件的支持。

> go get github.com/tnngo/log

## 示例

**快速使用：**

``` golang
log.NewSimple()
log.L().Debug("这是debug信息", zap.String("hello","golang"))
```

**日志文件：**
``` golang
opt := &Options{
	File: &FileCfg{
		Filename:   "test.log", // 文件名称。
		MaxSize:    512, // 最大尺寸MB，该示例表示512MB。
		MaxBackups: 100, // 最大备份数量，该示例表示100个。
		MaxAge:     30, // 最大保存时间，该示例表示30天。
		Compress:   true, // 是否压缩打包，该示例表示打包。
	},
}

log.New(opt)
// 使用默认打印方式，打印信息和参数信息都是零拷贝实现。
log.L().Info("这是debug信息", zap.String("hello","golang"))
```

## 默认日志级别

**控制台：** DEBUG。  
**文件：** INFO。
