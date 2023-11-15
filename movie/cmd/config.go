package main

type serviceConfig struct {
	GRPCConfig   grpcConfig   `yaml:"grpc"`
	ConsulConfig consulConfig `yaml:"consul"`
	HTTPConfig   httpConfig   `yaml:"http"`
}
type grpcConfig struct {
	Port string `yaml:"port"`
	Addr string `yaml:"addr"`
}
type httpConfig struct {
	Port string `yaml:"port"`
	Addr string `yaml:"addr"`
}
type consulConfig struct {
	Port string `yaml:"port"`
	Addr string `yaml:"addr"`
}
