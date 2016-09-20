package tools

type StreamConfig struct {
	// 流的名称
	StreamName string
	// 流的数据轮转方式
	TransitType int
	// 缓存区的大小
	CacheSize int
	// producers的配置
	ProducerConfigs []ProducerConfig
	// consumers的配置
	ConsumerConfigs []ConsumerConfig
}

type ProducerConfig struct {
	StreamName string
	// 当前producer的名字
	WorkerName string
	// current日志的位置
	FilePath string
	// 切割后的日志的位置
	BackDir string
	// 0表示从当前文件头开始读，1表示从backup文件夹的第一个文件开始读
	Position int
	// 日志切割的方式，通常是后缀为时间
	RollType int
	// producer的类型
	Producer string
	// 分隔符，这个应该是一个可选的配置，对于文本等数据默认以"\n"作为分隔符
	Delimiter string
	// 缓存的数据的大小
	BufSize int
	// 数据过滤器
	Filters []string
}

type ConsumerConfig struct {
	StreamName string
	WorkerName string
	Consumer   string
	// 输出的文件位置
	FilePath string
	///////////////////////////////////
	//  App模式，启动子程序处理数据  //
	///////////////////////////////////
	// 启动脚本的位置
	StartupScript string
	// 子进程的标准输出
	OutputDir string
}
