package parse

import "testing"

func TestLanguage_String(t *testing.T) {
	tests := []struct {
		name     string
		l        Language
		wantName string
	}{
		{name: "gb", l: GB, wantName: LanguageString[GB]},
		{name: "big", l: BIG5, wantName: LanguageString[BIG5]},
		{name: "gb|big", l: GB | BIG5, wantName: LanguageString[GB|BIG5]},
		{name: "gb|big|jp", l: GB | BIG5 | JP, wantName: LanguageString[GB|BIG5|JP]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotName := tt.l.String(); gotName != tt.wantName {
				t.Errorf("String() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
