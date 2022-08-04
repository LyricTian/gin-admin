package generate

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var commonInitialisms = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, cases.Title(language.English).String(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func ToLowerUnderlined(v string) string {
	if v == "" {
		return ""
	}

	var (
		value                                    = commonInitialismsReplacer.Replace(v)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase {
				if lastCase && (nextCase || nextNumber) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase && !nextNumber) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = true
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])
	s := strings.ToLower(buf.String())
	return s
}

func ToPlural(v string) string {
	return inflection.Plural(v)
}

func ToFirstLower(v string) string {
	return strings.ToLower(v[:1]) + v[1:]
}

func ToLower(v string) string {
	return strings.ToLower(v)
}

func ToLowerSpace(v string) string {
	return strings.Replace(ToLowerUnderlined(v), " ", "_", -1)
}

// --------------------------------- private methods ---------------------------------

const delimiter = "\n"

func execGoFmt(name string) error {
	cmd := exec.Command("gofmt", "-w", name, name)
	return cmd.Run()
}

func execParseTpl(tpl string, data interface{}) (*bytes.Buffer, error) {
	t := template.Must(template.New("").Parse(tpl))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func existsFile(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func writeFile(name string, buf *bytes.Buffer) error {
	_ = os.MkdirAll(filepath.Dir(name), os.ModePerm)
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, _ = io.Copy(file, buf)
	return nil
}

func readFile(name string) (*bytes.Buffer, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, file)
	return buf, nil
}

func insertContent(name string, fn func(string) (string, int, bool)) error {
	buf, err := readFile(name)
	if err != nil {
		return err
	}

	nbuf := new(bytes.Buffer)
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		cline := scanner.Text()
		data, flag, ok := fn(cline)
		if ok {
			if flag == -1 {
				nbuf.WriteString(data)
				nbuf.WriteString(delimiter)
				nbuf.WriteString(cline)
				nbuf.WriteString(delimiter)
				continue
			}
			nbuf.WriteString(cline)
			nbuf.WriteString(delimiter)
			nbuf.WriteString(data)
			nbuf.WriteString(delimiter)
			continue
		}
		nbuf.WriteString(cline)
		nbuf.WriteString(delimiter)
	}

	return writeFile(name, nbuf)
}

// --------------------------------- public methods ---------------------------------
