package main

import (
	"os"
	"testing"

	"github.com/tj/assert"
)

func TestEntitGen_buildTypeInfos(t *testing.T) {
	var g = &EntitGen{
		Output:       "./example/out",
		ProtobufType: "User",
		ProtoGenPath: "./example/gen/proto",
		ModelType:    "User",
		Options: Options{
			GenSlice: true,
			NoModels: []string{"Member"},
			NoEmbed:  []string{"Friends"},
			Suppress: []struct {
				Model string
				Field string
			}{
				{
					Field: "Permissions",
				},
			},
		},
		Debug: true,
	}
	g.init()

	pkgs, err := g.parsePackage(g.ProtoGenPath)
	assert.NoError(t, err)
	assert.NotNil(t, pkgs)

	outPkgs, err := g.parsePackage(g.Output)
	assert.NoError(t, err)
	assert.NotNil(t, outPkgs)

	typeinfos, err := g.buildTypeInfos(pkgs)
	assert.NoError(t, err)
	assert.NotNil(t, typeinfos)

	modelinfos, err := g.buildTypeInfos(outPkgs)
	assert.NoError(t, err)
	assert.NotNil(t, modelinfos)

	uinfo, ok := typeinfos["User"]
	assert.True(t, ok)
	assert.NotNil(t, uinfo)

	for _, f := range uinfo.Fields() {
		t.Logf("field %s type %v", f.Name(), f.Type())
	}

	msets := uinfo.Methods()
	assert.NotNil(t, msets)

	t.Logf("filename %s", uinfo.GetFilename())
	// for _, m := range uinfo.Methods() {
	// 	// t.Logf("method %s", m))
	// 	_ = m
	// }

	gens, err := g.convertPbToModel(uinfo, typeinfos, modelinfos, "User")
	assert.NoError(t, err)
	assert.NotNil(t, gens)

	t.Logf("genType %v", gens)
	t.Logf("generate output ----")
	for model, gentyp := range gens {
		t.Logf("model %s", model)
		g.writeTypeTo(gentyp, os.Stdout)
	}
}

func TestEntitGen_AliasName(t *testing.T) {
	var g = &EntitGen{
		Output:       "./example/out",
		ProtobufType: "User",
		ProtoGenPath: "./example/gen/proto",
		ModelType:    "UserI",
		Options: Options{
			GenSlice: true,
			NoModels: []string{"Member"},
			NoEmbed:  []string{"Friends"},
			Suppress: []struct {
				Model string
				Field string
			}{
				{
					Field: "Permissions",
				},
			},
		},
	}

	pkgs, err := g.parsePackage(g.ProtoGenPath)
	assert.NoError(t, err)
	assert.NotNil(t, pkgs)

	outPkgs, err := g.parsePackage(g.Output)
	assert.NoError(t, err)
	assert.NotNil(t, outPkgs)

	typeinfos, err := g.buildTypeInfos(pkgs)
	assert.NoError(t, err)
	assert.NotNil(t, typeinfos)

	modelinfos, err := g.buildTypeInfos(outPkgs)
	assert.NoError(t, err)
	assert.NotNil(t, modelinfos)

	uinfo, ok := typeinfos["User"]
	assert.True(t, ok)
	assert.NotNil(t, uinfo)

	for _, f := range uinfo.Fields() {
		t.Logf("field %s type %v", f.Name(), f.Type())
	}

	msets := uinfo.Methods()
	assert.NotNil(t, msets)

	t.Logf("filename %s", uinfo.GetFilename())
	// for _, m := range uinfo.Methods() {
	// 	// t.Logf("method %s", m))
	// 	_ = m
	// }

	gens, err := g.convertPbToModel(uinfo, typeinfos, modelinfos, "UserI")
	assert.NoError(t, err)
	assert.NotNil(t, gens)

	t.Logf("genType %v", gens)

	t.Logf("generate output ----")
	for model, gentyp := range gens {
		t.Logf("model %s", model)
		g.writeTypeTo(gentyp, os.Stdout)
	}
}

func TestEntitGen_buildType(t *testing.T) {
	var g = &EntitGen{
		Output:       "./example/out",
		ProtobufType: "User",
		ProtoGenPath: "./example/gen/proto",
		ModelType:    "User",
		Options: Options{
			GenSlice: true,
			NoModels: []string{"Member"},
			NoEmbed:  []string{"Friends"},
			Suppress: []struct {
				Model string
				Field string
			}{
				{
					Field: "Permissions",
				},
			},
		},
		Debug: true,
	}
	g.init()

	pkgs, err := g.parsePackage(g.ProtoGenPath)
	assert.NoError(t, err)
	assert.NotNil(t, pkgs)

	typeinfos, err := g.buildTypeInfos(pkgs)
	assert.NoError(t, err)
	assert.NotNil(t, typeinfos)
}

func TestEntitGen_Nested(t *testing.T) {
	var g = &EntitGen{
		Output:       "./example/out",
		ProtobufType: "Corp",
		ProtoGenPath: "./example/gen/proto",
		ModelType:    "Corp",
		Options: Options{
			GenSlice: true,
		},
		Debug: true,
	}
	g.init()

	pkgs, err := g.parsePackage(g.ProtoGenPath)
	assert.NoError(t, err)
	assert.NotNil(t, pkgs)

	outPkgs, err := g.parsePackage(g.Output)
	assert.NoError(t, err)
	assert.NotNil(t, outPkgs)

	typeinfos, err := g.buildTypeInfos(pkgs)
	assert.NoError(t, err)
	assert.NotNil(t, typeinfos)

	modelinfos, err := g.buildTypeInfos(outPkgs)
	assert.NoError(t, err)
	assert.NotNil(t, modelinfos)

	corpinfo, ok := typeinfos["Corp"]
	assert.True(t, ok)
	assert.NotNil(t, corpinfo)

	for _, f := range corpinfo.Fields() {
		t.Logf("field %s type %v", f.Name(), f.Type())
	}

	msets := corpinfo.Methods()
	assert.NotNil(t, msets)

	t.Logf("filename %s", corpinfo.GetFilename())
	// for _, m := range uinfo.Methods() {
	// 	// t.Logf("method %s", m))
	// 	_ = m
	// }

	gens, err := g.convertPbToModel(corpinfo, typeinfos, modelinfos, "Corp")
	assert.NoError(t, err)
	assert.NotNil(t, gens)

	t.Logf("genType %v", gens)
	t.Logf("generate output ----")
	for model, gentyp := range gens {
		t.Logf("model %s", model)
		g.writeTypeTo(gentyp, os.Stdout)
	}
}
