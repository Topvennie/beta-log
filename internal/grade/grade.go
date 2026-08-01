package grade

import (
	"fmt"
	"strconv"

	"github.com/Topvennie/beta-log/internal/database/model"
)

type Grade int

func (g Grade) Format(system model.GradeSystem) string {
	if g == 0 {
		return "?"
	}

	switch system {
	case model.GradeSystemFont:
		return g.font()
	case model.GradeSystemV:
		return g.v()
	default:
		return "n/a"
	}
}

func (g Grade) font() string {
	base := int(g) / 100
	remainder := int(g) % 100

	switch {
	case remainder == 0:
		return fmt.Sprintf("%da", base)
	case remainder <= 17:
		return fmt.Sprintf("%da+", base)
	case remainder <= 33:
		return fmt.Sprintf("%db", base)
	case remainder <= 50:
		return fmt.Sprintf("%db+", base)
	case remainder <= 67:
		return fmt.Sprintf("%dc", base)
	default:
		return fmt.Sprintf("%dc+", base)
	}
}

func (g Grade) v() string {
	switch {
	case g < 500:
		return "V0"
	case g < 550:
		return "V1"
	case g < 600:
		return "V2"
	default:
		font := g.font()
		if v, ok := fontToV[font]; ok {
			return v
		}

		return "n/a"
	}
}

// FromString converts a formatted grade to the internal grade
// It returns -1 if it is invalid
func FromString(s string, system model.GradeSystem) Grade {
	if s == "" {
		return 0
	}

	switch system {
	case model.GradeSystemFont:
		return fromFont(s)
	case model.GradeSystemV:
		return fromV(s)
	default:
		return -1
	}
}

func fromFont(s string) Grade {
	if len(s) < 2 || len(s) > 3 {
		return -1
	}

	base, err := strconv.Atoi(string(s[0]))
	if err != nil {
		return -1
	}

	suf := s[1:]
	var rem int

	switch suf {
	case "a":
		rem = 0
	case "a+":
		rem = 17
	case "b":
		rem = 33
	case "b+":
		rem = 50
	case "c":
		rem = 67
	case "c+":
		rem = 83
	default:
		return -1
	}

	return Grade(base*100 + rem)
}

func fromV(s string) Grade {
	switch s {
	case "V0":
		return 400
	case "V1":
		return 500
	case "V2":
		return 550
	default:
		font, ok := vToFont[s]
		if !ok {
			return -1
		}
		return fromFont(font)
	}
}

var fontToV = map[string]string{
	"6a":  "V3",
	"6a+": "V3",
	"6b":  "V4",
	"6b+": "V4",
	"6c":  "V5",
	"6c+": "V5",
	"7a":  "V6",
	"7a+": "V7",
	"7b":  "V8",
	"7b+": "V8",
	"7c":  "V9",
	"7c+": "V10",
	"8a":  "V11",
	"8a+": "V12",
	"8b":  "V13",
	"8b+": "V14",
	"8c":  "V15",
	"8c+": "V16",
	"9a":  "V17",
}

var vToFont = map[string]string{
	"V3":  "6a",
	"V4":  "6b",
	"V5":  "6c",
	"V6":  "7a",
	"V7":  "7a+",
	"V8":  "7b",
	"V9":  "7c",
	"V10": "7c+",
	"V11": "8a",
	"V12": "8a+",
	"V13": "8b",
	"V14": "8b+",
	"V15": "8c",
	"V16": "8c+",
	"V17": "9a",
}
