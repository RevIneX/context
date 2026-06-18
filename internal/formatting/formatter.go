package formatting

import (
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/RevIneX/context/internal/data"
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
	data.CatLanguage:     "header_how",
	data.CatRuntime:      "header_how",
	data.CatFramework:    "header_how",
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
	data.CatWebServer:     "header_where",
	data.CatPort:          "header_where",
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

	// Глобальные максимальные ширины
	maxLabelRunes := 0
	maxNameRunes := 0
	maxFileBytes := 0
	for _, f := range findings {
		label := categoryLabels[f.Category]
		if r := utf8.RuneCountInString(label); r > maxLabelRunes {
			maxLabelRunes = r
		}
		if r := utf8.RuneCountInString(f.Name); r > maxNameRunes {
			maxNameRunes = r
		}
		if len(f.File) > maxFileBytes {
			maxFileBytes = len(f.File)
		}
	}
	if maxLabelRunes < 6 {
		maxLabelRunes = 6
	}
	if maxNameRunes < 6 {
		maxNameRunes = 6
	}
	maxNameRunes += 1
	if maxFileBytes < 10 {
		maxFileBytes = 10
	}

	const foundInStr = "обнаружено в "

	groups := make(map[string][]data.Finding)
	for _, f := range findings {
		group, ok := categoryGroups[f.Category]
		if !ok {
			group = "header_inside"
		}
		groups[group] = append(groups[group], f)
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

		if gn.Key == "header_what" {
			wherePos := maxLabelRunes + 3 + maxNameRunes + 1 - 1
			numberPos := wherePos + len(foundInStr) + maxFileBytes - 8

			sb.WriteString("ЧТО ЭТО")
			for i := utf8.RuneCountInString("ЧТО ЭТО"); i < wherePos; i++ {
				sb.WriteByte(' ')
			}
			sb.WriteString("ГДЕ ЭТО")
			for i := wherePos + utf8.RuneCountInString("ГДЕ ЭТО"); i < numberPos; i++ {
				sb.WriteByte(' ')
			}
			sb.WriteString("НОМЕР СТРОКИ")
		} else {
			sb.WriteString(gn.Label)
		}
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
					writeRuneAligned(&sb, label, maxLabelRunes, f.Name, maxNameRunes, f.File, maxFileBytes, f.Line)
				} else {
					writeRuneAligned(&sb, "", maxLabelRunes, f.Name, maxNameRunes, f.File, maxFileBytes, f.Line)
				}
			}
		}
	}

	return sb.String()
}

func writeRuneAligned(sb *strings.Builder, label string, labelWidth int, name string, nameWidth int, file string, fileWidth int, line int) {
	labelRunes := utf8.RuneCountInString(label)
	sb.WriteString(label)
	for i := labelRunes; i < labelWidth; i++ {
		sb.WriteByte(' ')
	}
	sb.WriteString(" : ")

	nameRunes := utf8.RuneCountInString(name)
	sb.WriteString(name)
	for i := nameRunes; i < nameWidth; i++ {
		sb.WriteByte(' ')
	}

	sb.WriteString("обнаружено в ")

	sb.WriteString(file)
	for i := len(file); i < fileWidth; i++ {
		sb.WriteByte(' ')
	}

	sb.WriteString(" : ")
	sb.WriteString(strconv.Itoa(line))
	sb.WriteByte('\n')
}