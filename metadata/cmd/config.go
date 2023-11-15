package main

type serviceConfig struct {
	GRPCConfig   grpcConfig   `yaml:"grpc"`
	ConsulConfig consulConfig `yaml:"consul"`
	DBConfig     dbConfig     `yaml:"database"`
}
type grpcConfig struct {
	Port string `yaml:"port"`
	Addr string `yaml:"addr"`
}
type consulConfig struct {
	Port string `yaml:"port"`
	Addr string `yaml:"addr"`
}
type dbConfig struct {
	StrConn string `yaml:"strConn"`
}
