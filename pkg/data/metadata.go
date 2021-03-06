package data

// Metadata 描述资源相关的信息
type Metadata struct {
	Fansub []string
	Quality
	Type
	Language
	SubType
}

type (
	Type     uint8
	Quality  uint8
	Language uint8
	SubType  uint8
)

const (
	UnknownType Type = iota
	Full
	Episode
)

func (q Type) String() string {
	return CategoryString[q]
}

const (
	UnknownQuality Quality = iota
	P720
	P1080
	K2
)

func (q Quality) String() string {
	return QualityString[q]
}

const (
	UnknownSubType SubType = iota
	Internal
	External
)

func (s SubType) String() string {
	return SubTypeString[s]
}

const (
	UnknownLanguage          = 0
	GB              Language = 1 << iota
	BIG5
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
		GB:             "简体",
		BIG5:           "繁体",
		GB | BIG5:      "简繁",
		JP:             "日语",
		GB | JP:        "简日",
		BIG5 | JP:      "繁日",
		GB | BIG5 | JP: "简繁日",
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
	CategoryString = map[Type]string{
		UnknownType: "未知",
		Full:        "季度全集",
		Episode:     "单集",
	}
)
