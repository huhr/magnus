## Coding Robot has gone wild

## Version 1
想法有很多，第一个版本，想做到的就是可用，支持常用的producer和consumer，简单的负载均衡策略和
并发控制，梳理整体的模块逻辑架构，后续升级优化再逐步完善。

## 关于数据流
这里定义Stream就是业务上的一条数据流，例如业务模块需要收集的独立的日志流。一条数据流可能有多个
producer和consumer，其中producer负责生产数据，consumer负责消费数据，典型生产者消费者模型。每一
条Stream通过一个toml配置文件来描述，magnus启动时会检测conf下所有的配置文件，初始化并启动所有的
流。

Stream负责将producer生产的数据分发给consumer进行消费，最简单的方式可能是Stream按照顺序处理producer
生产的数据，实现一些简单的分发规则例如：consumer roundrobin, broadcast，带权重的随机分发等，以及
重试机制，可以理解为consumer消费是串行的逻辑。还需要考虑如何实现多consumers并行消费，以及相应的重
试机制。

producer需要支持多种类型，最简单的从console中获取数据，到从kafka、线上轮转日志等方式获取数据。producer
最重要的逻辑是对runtime信息的记录，保证程序可以正常的启停、异常崩溃时能够尽可能的不丢数据。现在已经
实现轮转日志实时收集的功能，这里需要简单的定义日志的轮转切割规则。

comsumer负责消费数据，可能的consumer包括归档数据到磁盘文本、数据发送给kafka、启动子进程并将数据发送给
子进程处理、发送给http service、发给送已经明确定义的thrift server等逻辑。已经实现发送子进程的逻辑。

magnus需要支持运行时对某一条Stream进行停止、重启等逻辑，需要通过signal实现，待完成。

magnus需要通过内置的thrift service将当前magnus的状态信息暴露出去，这些信息可能包括各个stream数据流转
量，异常记数，运行状态等等，通过另外的一套系统采集这些信息，并将信息汇总整理，通过一个统一的web管理
平台可以观察到magnus各项状态数据以及各个流的状态，以及配置一些监控项等等功能，待完成。


## 多producers、consumers的支持
1、并发控制    
2、是否保证顺序、是否保证不丢数据    

