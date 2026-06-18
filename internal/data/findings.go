package data

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
)

type Category string

const (
	CatArchitecture  Category = "architecture"
	CatLanguage      Category = "language"
	CatFramework     Category = "framework"
	CatDatabase      Category = "database"
	CatCache         Category = "cache"
	CatQueue         Category = "queue"
	CatOrchestration Category = "orchestration"
	CatContainer     Category = "container"
	CatMonitoring    Category = "monitoring"
	CatAPI           Category = "api"
	CatIaC           Category = "iac"
	CatCI            Category = "ci"
	CatCD            Category = "cd"
	CatEnv           Category = "env"
	CatWebServer     Category = "webserver"
	CatPort          Category = "port"
	CatClient        Category = "client"
	CatServer        Category = "server"
	CatRuntime       Category = "runtime"
	CatUnknown       Category = "unknown"
)

type Finding struct {
	Category Category          `json:"category"`
	Name     string            `json:"name"`
	Value    string            `json:"value,omitempty"`
	File     string            `json:"file"`
	Line     int               `json:"line"`
	Extra    map[string]string `json:"extra,omitempty"`
	Hash     string            `json:"hash"`
}

func (f *Finding) ComputeHash() {
	f.Hash = ""
	payload := f.serialize()
	h := sha256.Sum256([]byte(payload))
	f.Hash = hex.EncodeToString(h[:])
}

func (f *Finding) serialize() string {
	var sb strings.Builder
	sb.WriteString("category=")
	sb.WriteString(string(f.Category))
	sb.WriteByte('\n')
	sb.WriteString("name=")
	sb.WriteString(f.Name)
	sb.WriteByte('\n')
	sb.WriteString("value=")
	sb.WriteString(f.Value)
	sb.WriteByte('\n')
	sb.WriteString("file=")
	sb.WriteString(f.File)
	sb.WriteByte('\n')
	sb.WriteString("line=")
	sb.WriteString(strconv.Itoa(f.Line))
	sb.WriteByte('\n')

	keys := make([]string, 0, len(f.Extra))
	for k := range f.Extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString("extra.")
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(f.Extra[k])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (f *Finding) Validate() bool {
	if f.Hash == "" {
		return false
	}
	oldHash := f.Hash
	f.ComputeHash()
	newHash := f.Hash
	f.Hash = oldHash
	return oldHash == newHash
}

func NewFinding(cat Category, name, value, file string, line int, extra map[string]string) Finding {
	f := Finding{
		Category: cat,
		Name:     name,
		Value:    value,
		File:     file,
		Line:     line,
		Extra:    make(map[string]string, len(extra)),
	}
	for k, v := range extra {
		f.Extra[k] = v
	}
	f.ComputeHash()
	return f
}