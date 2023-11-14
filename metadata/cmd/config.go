package main

type serviceConfig struct {
	APIConfig apiConfig `yaml:"api"`
}
type apiConfig struct {
	Port       string `yaml:"port"`
	PortConsul string `yaml:"port_consul"`
}
