package queue

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	}
	Log struct {
		Path  string `yaml:"path"`
		File  string `yaml:"file"`
		Level string `yaml:"level"`
	}
	Queue struct {
		Handle_chan_size int `yaml:"handle_chan_size"`
		Max_wait_count   int `yaml:"max_wait_count"`
	}
}
