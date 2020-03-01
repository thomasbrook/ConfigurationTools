# Configuration-tools
基于go语言，UI类库walk，实现windows 桌面版配置工具

# GUI walk
https://github.com/lxn/walk

# 配置文件工具
go get github.com/akavel/rsrc
rsrc -manifest test.manifest -o rsrc.syso

# 从 io.Reader 中读数据
https://cloud.tencent.com/developer/article/1422483

# 隐藏cmd窗口
go build -ldflags="-H windowsgui"
