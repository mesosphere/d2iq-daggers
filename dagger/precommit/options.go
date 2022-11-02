package precommit

type config struct {
	baseImage string
}

func defaultConfig() config {
	return config{
		baseImage: "python:3.12.0a1-bullseye",
	}
}

type Option func(config) config

func BaseImage(img string) Option {
	return func(c config) config {
		c.baseImage = img
		return c
	}
}
