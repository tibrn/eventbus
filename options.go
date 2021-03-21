package eventbus

import "runtime"

type Options struct {
	Concurrency       int
	RunnerConcurrency int
}

func defaultOrOptions(options ...*Options) *Options {

	finalOptions := &Options{
		Concurrency:       runtime.GOMAXPROCS(0) * 4,
		RunnerConcurrency: runtime.GOMAXPROCS(0) * 4,
	}

	for _, opt := range options {
		if opt.Concurrency > 0 {
			finalOptions.Concurrency = opt.Concurrency
		}
		if opt.RunnerConcurrency > 0 {
			finalOptions.RunnerConcurrency = opt.RunnerConcurrency
		}
	}

	return finalOptions
}
