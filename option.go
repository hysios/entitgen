package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Options struct {
	NoModels []string
	Suppress []struct {
		Model string
		Field string
	}
	GenSlice         bool
	EntityFileDetect bool
	NoEmbed          []string
	sets             map[string]SetOptionFunc
}

// parseOptions
func parseOptions(options []string) (*Options, error) {
	var (
		outOpts = Options{}
	)
	if err := outOpts.Apply(options); err != nil {
		return nil, err
	}

	outOpts.GenSlice = genSlice
	return &outOpts, nil
}

type SetOptionFunc func(opts *Options, value interface{}) error

// AddMethod
func (opts *Options) AddMethod(key string, set SetOptionFunc) {
	if opts.sets == nil {
		opts.sets = make(map[string]SetOptionFunc)
	}

	opts.sets[key] = set
}

// Apply
func (opts *Options) Apply(ops []string) error {
	if opts.sets == nil {
		opts.sets = GlobalOption.sets
	}

	for _, op := range ops {
		es := strings.Split(op, "=")
		key, value := es[0], es[1]

		if set, ok := opts.sets[key]; ok {
			if err := set(opts, value); err != nil {
				return err
			}
		}
	}
	return nil
}

var GlobalOption = Options{}

func init() {
	GlobalOption.AddMethod("no-models", func(opts *Options, value interface{}) error {
		opts.NoModels = append(opts.NoModels, value.(string))
		return nil
	})

	GlobalOption.AddMethod("suppress", func(opts *Options, value interface{}) error {
		lines := strings.Split(value.(string), ",")
		for _, v := range lines {
			es := strings.Split(v, ".")
			if len(es) != 2 {
				opts.Suppress = append(opts.Suppress, struct {
					Model string
					Field string
				}{Field: es[0]})
			} else {

				opts.Suppress = append(opts.Suppress, struct {
					Model string
					Field string
				}{Model: es[0], Field: es[1]})
			}
		}
		return nil
	})

	GlobalOption.AddMethod("entity-file-detect", func(opts *Options, value interface{}) error {
		v, err := strconv.ParseBool(fmt.Sprint(value))
		if err != nil {
			return err
		}
		opts.EntityFileDetect = v
		return nil
	})

	GlobalOption.AddMethod("no-embed", func(opts *Options, value interface{}) error {
		lines := strings.Split(value.(string), ",")
		for _, v := range lines {
			opts.NoEmbed = append(opts.NoEmbed, v)
		}
		return nil
	})

}
