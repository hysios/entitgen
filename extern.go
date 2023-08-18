package main

import (
	"strings"

	"github.com/hysios/entitgen/gen"
)

type ExternalType struct {
	Imports []gen.Pkg
	Type    gen.Var
	Convert gen.Converter
}

var (
	// ExternalTypes
	externalTypes = map[string]ExternalType{
		"google.golang.org/protobuf/types/known/timestamppb.Timestamp": ExternalType{
			Imports: []gen.Pkg{
				gen.Pkg{
					Fullname: "google.golang.org/protobuf/types/known/timestamppb",
				},
			},
			Type: gen.Var{
				Name:    "t",
				Type:    "Time",
				PkgName: "time",
			},
			Convert: &externConvert{
				toProto: func(in string) string {
					return in + ".AsTime()"
				},
				fromProto: func(in string) string {
					return "timestamppb.New(" + in + ")"
				},
			},
		},
	}
)

type externConvert struct {
	fromProto func(in string) string
	toProto   func(in string) string
}

// From
func (c *externConvert) From(in string) string {
	return c.fromProto(in)
}

// To
func (c *externConvert) To(in string) string {
	return c.toProto(in)
}

func isExternalType(typ string) bool {
	// if typ has star
	typ = strings.TrimPrefix(typ, "*")
	if _, ok := externalTypes[typ]; ok {
		return true
	}
	return false
}

func getExternalType(typ string) (ExternalType, bool) {
	typ = strings.TrimPrefix(typ, "*")
	t, ok := externalTypes[typ]
	return t, ok
}

// getExternalConvert
func getExternalConvert(typ string) gen.Converter {
	ext, ok := getExternalType(typ)
	if !ok {
		return nil
	}

	return ext.Convert
}
