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
	WorkerName string
	FilePath   string
	BackDir    string
	Position   int
	RollType   int
	Producer   string
	Delimiter  string
	BufSize    int
	Filters    []string
}

type ConsumerConfig struct {
	StreamName string
	WorkerName string
	Consumer   string
	FilePath   string
	Filters    []string
	///////////////////////////////////
	//  App模式，启动子程序处理数据  //
	///////////////////////////////////
	StartupScript string
	OutputDir     string
}
