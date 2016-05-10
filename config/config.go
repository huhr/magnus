package config

type StreamConfig struct {
	Name string
	RollType int
	CacheSize int
	Pcfgs []ProducerConfig
	Ccfgs []ConsumerConfig
}

type ProducerConfig struct {
	Producer string
}

type ConsumerConfig struct {
	Consumer string
}
