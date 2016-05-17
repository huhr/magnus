### 大体逻辑
程序的核心是一个Adapter实例，该实例读取完配置文件后，初始化producers，pipeline，filters，consumers。
producers获取数据，经过filter，传输到pipeline中，这里pipeline是一组有一定缓存的channel，每一个数据流
对应一个pipeline，类似kafka的topic的概念，pipeline两端可能是多个producers或者consumers。    

### 主要的实现功能和模式
提供n:m的功能很重要，例如logtailer，如果后端python的app，出现性能问题时，
依赖kafka的partitions才能水平扩展，如果上游不支持水平扩展，比较麻烦，支持
类似actsender2进行的水平扩展很关键，支持简单轮询，权重等分配规则。

### 解决哪些问题
构想：支持支持扩展开发模式，类似actsender2的场景使用开发模式，定制producer
重新编译，consumer使用多个app即可。类似日志收集等需求，adapter直接提供。
改善jvm资源消耗，启动慢等问题。

### 大体逻辑
producer、filter、format、dispatcher、format、filter、consumer
