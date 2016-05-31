## Coding Robot has gone wild

tail文件    
touch text    
echo xxx.log >> ./text    
make    
./magnus -d `pwd`    


## Magnus
包含三个模块:
	部署在集群机器上，负责进行数据收集转发的agent实例，也就是adapter模块
	负责配置启动机器上adapter，展示adapter状态及各个数据流状态的web模块
	负责对集群adapter进行状态数据收集存储，已经进行监控报警的Monitor模块

