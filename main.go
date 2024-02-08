package main

import (
	"flag"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	debug            bool
	genPath          string
	pkgName          string
	protobufPath     string
	protobufType     string
	protoGenPath     string
	modelType        string
	composeTypes     arrayFlags
	filenameTemplate string
	options          arrayFlags
	genSlice         bool
)

var log = zap.S()

func initLog(lvl int) *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zapcore.Level(lvl))
	log, _ := cfg.Build()
	return log.Sugar()
}

func main() {
	parseArgs()

	opts, err := parseOptions(options)
	if err != nil {
		log.Fatalf("parse options failed: %v", err)
	}
	log.Debugf("Options: %v", opts)
	var gen = EntitGen{
		PkgName:      pkgName,
		Output:       genPath,
		ProtobufType: protobufType,
		ProtoPath:    protobufPath,
		ProtoGenPath: protoGenPath,
		ModelType:    modelType,
		Options:      *opts,
		Debug:        debug,
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
	flag.BoolVar(&genSlice, "slice", false, "generate slice type")
	flag.Var(&options, "option", "option")
	flag.Var(&options, "O", "option")
	flag.BoolVar(&debug, "debug", true, "debug level")

	flag.Parse()

	root := lookupProjectDir()
	protoGenPath = resolvePath(root, protoGenPath)
	genPath = resolvePath(root, genPath)
}

func resolvePath(root string, rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(root, rel)
}

func lookupProjectDir() string {
	cwd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}

		cwd = filepath.Dir(cwd)
		if cwd == "/" {
			break
		}
	}

	return ""
}
