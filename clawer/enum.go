package clawer

type (
	Category uint8
	Quality  uint8
	Language uint8
	SubType  uint8
)

const (
	Unknown Category = iota
	FullSession
	Animate
	Movie
)

func (q Category) String() string {
	return CategoryString[q]
}

const (
	P720 Quality = iota
	P1080
	K2
)

func (q Quality) String() string {
	return QualityString[q]
}

const (
	Internal SubType = iota
	External
)

func (s SubType) String() string {
	return SubTypeString[s]
}

const (
	GB Language = 1 << iota
	BIG
	JP
)

func (l Language) String() (name string) {
	for lan, n := range LanguageString {
		if l == lan {
			return n
		}
		if l&lan == 1 {
			name += n
		}
	}
	return name
}

var (
	LanguageString = map[Language]string{
		GB:  "简体",
		BIG: "繁体",
		JP:  "日语",
	}
	QualityString = map[Quality]string{
		P1080: "1080p",
		P720:  "720p",
		K2:    "2k",
	}
	SubTypeString = map[SubType]string{
		Internal: "内置",
		External: "外挂",
	}
	CategoryString = map[Category]string{
		Unknown:     "未知",
		FullSession: "季度全集",
		Animate:     "动画",
		Movie:       "电影",
	}
)
