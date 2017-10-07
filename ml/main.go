package main

// build vars
var (
	Version string
	Build   string
	mlCli   = &mlCLI{}
	config  = &CliConfig{}
)

func main() {
	config.init(Version, Build)
	cli()
}
