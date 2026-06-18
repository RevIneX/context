package formatting

import (
	"context/internal/data"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var groupNames = []struct {
	Key   string
	Label string
}{
	{"header_what", "ЧТО ЭТО"},
	{"header_how", "НА ЧЁМ НАПИСАНО"},
	{"header_inside", "ЧТО ТУТ ЕСТЬ"},
	{"header_where", "ГДЕ ЧТО ХОСТИТСЯ"},
}

var categoryGroups = map[data.Category]string{
	data.CatArchitecture: "header_what",
	data.CatClient:       "header_what",
	data.CatServer:       "header_what",

	data.CatLanguage:  "header_how",
	data.CatRuntime:   "header_how",
	data.CatFramework: "header_how",

	data.CatEnv:           "header_inside",
	data.CatOrchestration: "header_inside",
	data.CatDatabase:      "header_inside",
	data.CatContainer:     "header_inside",
	data.CatMonitoring:    "header_inside",
	data.CatQueue:         "header_inside",
	data.CatCache:         "header_inside",
	data.CatAPI:           "header_inside",
	data.CatIaC:           "header_inside",
	data.CatCI:            "header_inside",
	data.CatCD:            "header_inside",
	data.CatUnknown:       "header_inside",

	data.CatWebServer: "header_where",
	data.CatPort:      "header_where",
}

var categoryLabels = map[data.Category]string{
	data.CatArchitecture:  "архитектура",
	data.CatClient:        "клиент",
	data.CatServer:        "сервер",
	data.CatLanguage:      "языки программирования",
	data.CatRuntime:       "среда выполнения",
	data.CatFramework:     "фреймворки",
	data.CatEnv:           "переменные окружения",
	data.CatOrchestration: "оркестрация",
	data.CatDatabase:      "базы данных",
	data.CatContainer:     "контейнеры",
	data.CatMonitoring:    "мониторинг",
	data.CatQueue:         "очереди",
	data.CatCache:         "кэши",
	data.CatAPI:           "API",
	data.CatIaC:           "IaC",
	data.CatCI:            "CI",
	data.CatCD:            "CD",
	data.CatUnknown:       "неопознанное",
	data.CatWebServer:     "веб-серверы",
	data.CatPort:          "порты",
}

func Format(findings []data.Finding) string {
	if len(findings) == 0 {
		return ""
	}

	groups := make(map[string][]data.Finding)
	for _, f := range findings {
		group, ok := categoryGroups[f.Category]
		if !ok {
			group = "header_inside"
		}
		groups[group] = append(groups[group], f)
	}

	maxLabelRunes := 0
	maxNameLen := 0
	maxFileLen := 0
	for cat := range categoryLabels {
		label := categoryLabels[cat]
		if r := utf8.RuneCountInString(label); r > maxLabelRunes {
			maxLabelRunes = r
		}
	}
	for _, f := range findings {
		if len(f.Name) > maxNameLen {
			maxNameLen = len(f.Name)
		}
		if len(f.File) > maxFileLen {
			maxFileLen = len(f.File)
		}
	}
	if maxNameLen < 18 {
		maxNameLen = 18
	}
	if maxFileLen < 32 {
		maxFileLen = 32
	}

	var sb strings.Builder
	firstGroup := true

	for _, gn := range groupNames {
		items, ok := groups[gn.Key]
		if !ok || len(items) == 0 {
			continue
		}

		if !firstGroup {
			sb.WriteByte('\n')
		}
		firstGroup = false

		sb.WriteString(gn.Label)
		sb.WriteByte('\n')

		byCategory := make(map[data.Category][]data.Finding)
		for _, f := range items {
			byCategory[f.Category] = append(byCategory[f.Category], f)
		}

		var cats []data.Category
		for cat := range byCategory {
			cats = append(cats, cat)
		}
		sort.Slice(cats, func(i, j int) bool {
			return string(cats[i]) < string(cats[j])
		})

		for _, cat := range cats {
			label := categoryLabels[cat]
			ff := byCategory[cat]

			seen := make(map[string]bool)
			var unique []data.Finding
			for _, f := range ff {
				key := f.Name + "|" + f.Value
				if !seen[key] {
					seen[key] = true
					unique = append(unique, f)
				}
			}
			sort.Slice(unique, func(i, j int) bool {
				return unique[i].Name < unique[j].Name
			})

			for i, f := range unique {
				if i == 0 {
					writeAligned(&sb, label, maxLabelRunes, f.Name, maxNameLen, f.File, maxFileLen, f.Line)
				} else {
					writeAligned(&sb, "", maxLabelRunes, f.Name, maxNameLen, f.File, maxFileLen, f.Line)
				}
			}
		}
	}

	return sb.String()
}

func writeAligned(sb *strings.Builder, label string, labelWidth int, name string, nameWidth int, file string, fileWidth int, line int) {
	labelRunes := utf8.RuneCountInString(label)
	sb.WriteString(label)
	for i := labelRunes; i < labelWidth+2; i++ {
		sb.WriteByte(' ')
	}
	sb.WriteString(" : ")

	sb.WriteString(name)
	for i := len(name); i < nameWidth; i++ {
		sb.WriteByte(' ')
	}

	sb.WriteString(" обнаружено в ")

	sb.WriteString(file)
	for i := len(file); i < fileWidth; i++ {
		sb.WriteByte(' ')
	}

	sb.WriteString(" : ")
	sb.WriteString(strconv.Itoa(line))
	sb.WriteByte('\n')
}