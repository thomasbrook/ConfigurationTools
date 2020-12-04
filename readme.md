#### 一、GUI walk
https://github.com/lxn/walk

#### 二、配置文件工具
go get github.com/akavel/rsrc
rsrc -manifest test.manifest -o rsrc.syso

#### 三、从 io.Reader 中读数据
https://cloud.tencent.com/developer/article/1422483

#### 四、隐藏cmd窗口
go build -ldflags="-H windowsgui"

#### 五、如何打包
将img/*、.ini(可无，会自动生成)、config.json、favicon.ico、help.html、logo.png、exe执行文件（第四步生成的exe文件）放置于自定义文件夹内。
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo16.png)
将这些文件添加到rar包内，选择可执行文件。具体打包方式请百度。

#### 六、主要功能
###### #1 首页
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo1.png)
###### #2 批量编辑CAN信息
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo6.png)
###### #3 从其他车型导入配置信息
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo8.png)
###### #4 编辑车系，可编辑分组信息以及设置是否为智能机等操作
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo4.png)
###### #5 创建车系
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo5.png)
###### #6 某车系下的CAN管理界面，可导出CSV或者复制到剪贴板
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo2.png)
###### #7 CAN管理界面中，双击行，可进行编辑
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo3.png)
###### #8 从剪贴板，导入配置信息
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo9.png)
###### #9 从CSV文件导入
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo10.png)
###### #10 从大数据导入，根据can关键词或者唯一编码搜索添加
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo7.png)
###### #11 导出设备的历史报文，可单台设备或者打开txt文件（每行一个设备编号）
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo11.png)
###### #12 多字段组合多枚举值，将多个字段进行组合展示
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo12.png)
###### #13 指令管理与发送（正则表达式输入校验，命令控件组合与参数动态拼接）
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo14.png)
![image](https://raw.githubusercontent.com/thomasbrook/ConfigurationTools/master-v2/img/demo15.png)
