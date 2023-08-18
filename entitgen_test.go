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
			NoModels: []string{"Member"},
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

	typeinfos, err := g.buildTypeInfos(pkgs)
	assert.NoError(t, err)
	assert.NotNil(t, typeinfos)

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

	gens, err := g.convertPbToModel(uinfo, typeinfos, nil, "User")
	assert.NoError(t, err)
	assert.NotNil(t, gens)

	t.Logf("genType %v", gens)

	t.Logf("generate output ----")
	for model, gentyp := range gens {
		t.Logf("model %s", model)
		g.writeTypeTo(gentyp, os.Stdout)
	}
}
