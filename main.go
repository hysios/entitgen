package main

import (
	"flag"
	"log"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myFlags arrayFlags
var (
	genPath          string
	pkgName          string
	protobufPath     string
	protobufType     string
	protoGenPath     string
	modelType        string
	composeTypes     arrayFlags
	filenameTemplate string
	options          arrayFlags
)

func main() {
	parseArgs()

	opts, err := parseOptions(options)
	if err != nil {
		log.Fatalf("parse options failed: %v", err)
	}
	log.Printf("Options: %v", opts)
	var gen = EntitGen{
		PkgName:      pkgName,
		Output:       genPath,
		ProtobufType: protobufType,
		ProtoPath:    protobufPath,
		ProtoGenPath: protoGenPath,
		ModelType:    modelType,
		Options:      *opts,
	}

	if err := gen.Run(); err != nil {
		log.Fatalf("running entitgen failed: %v", err)
	}
}

func parseArgs() {
	flag.StringVar(&genPath, "gen", "./model", "Path to the model directory.")
	flag.StringVar(&pkgName, "pkg", "model", "Path to the model directory.")
	// flag.StringVar(&protobufPath, "proto-path", "./proto", "Path to the protobuf directory.")
	flag.StringVar(&protoGenPath, "proto-gen-path", "./gen/proto", "Path of the protobuf gen directory")
	flag.StringVar(&protoGenPath, "G", "./gen/proto", "Path of the protobuf gen directory")
	flag.StringVar(&protobufType, "type", "", "generate entity Model form protobuf type")
	flag.StringVar(&protobufType, "T", "", "generate entity Model form protobuf type")
	flag.StringVar(&modelType, "model", "", "generate protobuf type to entity Model")
	flag.StringVar(&modelType, "M", "", "generate protobuf type to entity Model")
	flag.Var(&composeTypes, "compose-type", "compose many type to entity Model")
	flag.StringVar(&filenameTemplate, "filename-template", "{{.Name}}.entity.go", "filename template")
	flag.Var(&options, "option", "option")
	flag.Var(&options, "O", "option")

	flag.Parse()
}
