package config

type StreamConfig struct {
	Name		string
	TransitType	int
	CacheSize	int
	Pcfgs		[]ProducerConfig
	Ccfgs		[]ConsumerConfig
}

type ProducerConfig struct {
	StreamName  string
	WorkerName        string
	FilePath    string
	BackDir     string
	RollType    int
	Producer	string
	Delimiter	string
	BufSize		int
	Filters		[]string
}

type ConsumerConfig struct {
	StreamName  string
	WorkerName        string
	Consumer	string
	Filters		[]string
}

