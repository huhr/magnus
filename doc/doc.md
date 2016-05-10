### 大体逻辑
程序的核心是一个Adapter实例，该实例读取完配置文件后，初始化producers，pipeline，filters，consumers。
producers获取数据，经过filter，传输到pipeline中，这里pipeline是一组有一定缓存的channel，每一个数据流
对应一个pipeline，类似kafka的topic的概念，pipeline两端可能是多个producers或者consumers。    

