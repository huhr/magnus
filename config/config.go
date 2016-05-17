package config

type StreamConfig struct {
	Name		string
	RollType	int
	CacheSize	int
	Pcfgs		[]ProducerConfig
	Ccfgs		[]ConsumerConfig
}

type ProducerConfig struct {
	FilePath    string
	Producer	string
	Delimiter	string
	BufSize		int
	Filters		[]string
}

type ConsumerConfig struct {
	Consumer	string
	Filters		[]string
}

