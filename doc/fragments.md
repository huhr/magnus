### LICENSE
开源项目的第一个概念应该是License，比如Gollum使用的是Apache License 2.0，
而Heka使用Mozilla Public License 2.0，其他主流的License还有BSD Licenses，
GNU General Public License，MIT License等。    
参见http://www.awflasher.com/blog/archives/939

### vendor
Go程序通常由许多不同的package构成，这些package的源码来自标准库或者GOPATH
里的第三方库，我们通过修改GOPATH并将第三方库copy到项目里进行构建，这样Go
模块可能包含较多冗余第三方源码，并且还需要写install脚本。我们也可以统一管
理第三方里依赖，这样就会带来各种版本升级管理问题。    
官方在1.5版本中给出了一种解决方案，在保持原有功能不变的前提下，引入vendor
文件夹的概念。    
https://medium.com/@freeformz/go-1-5-s-vendor-experiment-fd3e830f52c3    
https://blog.gopheracademy.com/advent-2015/vendor-folder

### 配置文件
在simplelog模块中，使用json格式的配置文件，但是手写json是比较容易出错的，
可读性也一般，无法注释，优点是解析起来比较简单，使用标准库即可。相当来说
INI文件更好一些，但是这里使用TOML作为配置文件格式，TOML目前0.4.0版本仍不
是稳定版本，使用的toml解析库完整支持0.2.0版本规范。    
https://github.com/toml-lang/toml    
https://github.com/BurntSushi/toml

### Offset
Offset是以byte为单位的
