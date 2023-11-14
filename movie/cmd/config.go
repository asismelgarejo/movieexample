package main

type serviceConfig struct {
	APIConfig    apiConfig    `yaml:"api"`
	GRPCConfig   grpcConfig   `yaml:"grpc"`
	ConsulConfig consulConfig `yaml:"consul"`
}
type apiConfig struct {
	Port       string `yaml:"port"`
	PortConsul string `yaml:"port_consul"`
}
type grpcConfig struct {
	Port string `yaml:"port"`
}
type consulConfig struct {
	Port string `yaml:"port"`
}
