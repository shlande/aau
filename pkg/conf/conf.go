package conf

type Config struct {
	*Http
	*Log
	*Data
}

type Http struct {
	Address string
}

type Data struct {
	OutputDir string
}

type Log struct {
	Level  string
	Output string
}

var Default = Config{
	Http: &Http{
		Address: ":9099",
	},
	Log: &Log{
		Level:  "info",
		Output: "./subs.log",
	},
	Data: &Data{
		OutputDir: ".",
	},
}
