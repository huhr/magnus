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

### 文件系统
文件系统是对存储设备上的数据进行组织的机制。在Linux中讲一个文件系统与一个设    
备关联起来的过程成为挂装（mount）。文件系统类型如Ext2，Ext3，Ext4。 虚拟文件
系统位于文件系统的上层，隐藏不通的文件系统，对上层提供一个通用的文件模型。
Linux内核不对一个特定的函数进行硬编码来执行诸如read()这样的操作，而是对每个
操作都必须使用一个指针，指向到具体的文件系统的相应的函数。这个通用的文件模型
包含：superblock（存放已安装文件系统的有关信息）、inode（存放关于具体文件的
一般信息）、file（存放打开文件与进程之间交互的有关信息）、dentry（存放目录项
与对应的文件进行连接的有关信息）。

https://www.ibm.com/developerworks/cn/linux/l-linux-filesystem/figure1.gif


### nil != nil
Golang的interface类型包含两个元素，类型与值，仅当类型和值都为nil时，该interface
对象才为nil，例如

	var a *int
	a = nil 
	interface{}(a) != nil

### 文件编辑模式
修改文本时有通常有两种模式：一种是直接修改文件，不需要额外的磁盘空间，缺点是编
辑直接落地了，编辑过程中程序崩溃，文件数据就是不完整的了。另一种模式是生成一个
新的文件，完成编辑时进行mv操作，这需要更多的磁盘空间，优点是mv是原子操作，程序
即使崩溃，老的文件也没有被破坏。

### 匿名函数

	var nums = []int{1, 2, 3, 4, 5, 6}
	for _, num := range nums {
		go func() {
			fmt.Println(num)
		}()
	}
这里输出6个6，首先golang的loop，这里的num只声明了一次，后续的loop中其实是共用该    
变量，所以loop里的协程其实引用的同一个变量，所以会输出7个6。在这个场景里需要使用    
传参的方式来解决。

### 启动子进程
golang提供的os.StartProcess()通过系统调用SYS_CLONE来创建新的进程，使用EXEC执行，
使用SYS_DUP进行输入输出重定向
