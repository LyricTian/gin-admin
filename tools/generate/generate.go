package generate

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/util/json"
	"github.com/LyricTian/gin-admin/v9/pkg/util/yaml"
	"go.uber.org/zap"
)

type Config struct {
	Directory          string  `yaml:"-"`
	TplDirectory       string  `yaml:"-"`
	PkgName            string  `yaml:"pkg_name"`
	ModuleName         string  `yaml:"module_name"`
	ModuleLowerName    string  `yaml:"-"`
	Name               string  `yaml:"name"`
	PluralName         string  `yaml:"-"`
	LowerUnderlineName string  `yaml:"-"`
	LowerPluralName    string  `yaml:"-"`
	LowerSpaceName     string  `yaml:"-"`
	LowerName          string  `yaml:"-"`
	Description        string  `yaml:"description"`
	Fields             []Field `yaml:"fields"`
}

type Field struct {
	Name                string `yaml:"name"`
	LowerUnderlineName  string `yaml:"-"`
	FirstLowerName      string `yaml:"-"`
	Type                string `yaml:"type"`
	Optional            bool   `yaml:"optional"`
	GormTag             string `yaml:"gorm_tag"`
	BindingTag          string `yaml:"binding_tag"`
	Comments            string `yaml:"comments"`
	InQuery             bool   `yaml:"in_query"`
	QueryComments       string `yaml:"query_comments"`
	QueryIfExpression   string `yaml:"query_if_expression"`   // For example: `v := params.Username; v != ""`
	QueryWhereCondition string `yaml:"query_where_condition"` // For example: `"username=?", v`
	InCreate            bool   `yaml:"in_create"`
}

// Generate module files by template
func Generate(ctx context.Context, dir, tplDir string, results []string, yamlCfgFile string) error {
	buf, err := ioutil.ReadFile(yamlCfgFile)
	if err != nil {
		logger.Context(ctx).Error("Failed to read yaml file", zap.String("name", yamlCfgFile), zap.Error(err))
		return err
	}

	cfg := new(Config)
	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		return err
	}
	cfg.Directory = dir
	cfg.TplDirectory = tplDir
	cfg.ModuleLowerName = ToLower(cfg.ModuleName)
	cfg.LowerUnderlineName = ToLowerUnderlined(cfg.Name)
	cfg.PluralName = ToPlural(cfg.Name)
	cfg.LowerPluralName = ToLower(cfg.PluralName)
	cfg.LowerSpaceName = ToLowerSpace(cfg.Name)
	cfg.LowerName = ToLower(cfg.Name)

	for i, field := range cfg.Fields {
		cfg.Fields[i].LowerUnderlineName = ToLowerUnderlined(field.Name)
		cfg.Fields[i].FirstLowerName = ToFirstLower(field.Name)
	}

	fmt.Println("Use configuraton to generate code:")
	fmt.Println()
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonCfg))
	fmt.Println()

	for _, result := range results {
		if err := parseTplAndWriteFile(ctx, cfg, result); err != nil {
			return err
		}
	}

	moduleFile := filepath.Join(cfg.Directory, "internal", "module", cfg.ModuleLowerName, cfg.ModuleLowerName+".go")
	if !existsFile(moduleFile) {
		err := parseTplAndWriteFile(ctx, cfg, "module")
		if err != nil {
			return err
		}

		// Insert module to internal/init.go
		initFile := filepath.Join(cfg.Directory, "internal", "init.go")
		initModuleStart, initRegisterStart := false, false
		err = insertContent(initFile, func(line string) (string, int, bool) {
			line = strings.TrimSpace(line)
			if !initModuleStart && strings.HasSuffix(line, "// Initialize modules") {
				initModuleStart = true
			}
			if initModuleStart && strings.HasSuffix(line, "// end") {
				initModuleStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("if err := injector.%s.Init(ctx); err != nil {\n", cfg.ModuleName))
				buf.WriteString("return cleanInjectFn, err\n}")
				return buf.String(), -1, true
			}

			if !initRegisterStart && strings.HasPrefix(line, "apiGroup :=") {
				initRegisterStart = true
			}
			if initRegisterStart && strings.HasSuffix(line, "// end") {
				initRegisterStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("\n// Register %s APIs\n", cfg.ModuleName))
				buf.WriteString(fmt.Sprintf("injector.%s.RegisterAPI(ctx, apiGroup)", cfg.ModuleName))
				return buf.String(), -1, true
			}

			return "", 0, false
		})
		if err != nil {
			return err
		}

		// Insert module to internal/inject/wire.go
		wireFile := filepath.Join(cfg.Directory, "internal", "inject", "wire.go")
		wireImportStart, wireBuildStart := false, false
		err = insertContent(wireFile, func(line string) (string, int, bool) {
			line = strings.TrimSpace(line)
			if !wireImportStart && strings.HasPrefix(line, "import (") {
				wireImportStart = true
			}
			if wireImportStart && strings.HasSuffix(line, "// end") {
				wireImportStart = false
				var buf bytes.Buffer
				buf.WriteByte('\t')
				buf.WriteString(fmt.Sprintf(`"%s/internal/module/%s"`, cfg.PkgName, cfg.ModuleLowerName))
				return buf.String(), -1, true
			}

			if !wireBuildStart && strings.HasPrefix(line, "wire.Build") {
				wireBuildStart = true
			}
			if wireBuildStart && strings.HasSuffix(line, "// end") {
				wireBuildStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("\t\t%s.Set,", cfg.ModuleLowerName))
				return buf.String(), -1, true
			}
			return "", 0, false
		})
		if err != nil {
			return err
		}

		// Insert module to internal/inject/inject.go
		injectFile := filepath.Join(cfg.Directory, "internal", "inject", "inject.go")
		injectImportStart, injectStructStart := false, false
		err = insertContent(injectFile, func(line string) (string, int, bool) {
			line = strings.TrimSpace(line)
			if !injectImportStart && strings.HasPrefix(line, "import (") {
				injectImportStart = true
			}
			if injectImportStart && strings.HasSuffix(line, "// end") {
				injectImportStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf(`"%s/internal/module/%s"`, cfg.PkgName, cfg.ModuleLowerName))
				return buf.String(), -1, true
			}

			if !injectStructStart && strings.HasPrefix(line, "type Injector struct") {
				injectStructStart = true
			}
			if injectStructStart && strings.HasSuffix(line, "// end") {
				injectStructStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("%s  *%s.%s", cfg.ModuleName, cfg.ModuleLowerName, cfg.ModuleName))
				return buf.String(), -1, true
			}
			return "", 0, false
		})
		if err != nil {
			return err
		}
	} else {
		moduleSetStart, moduleStructStart, moduleMigrateStart, moduleGroupStart := false, false, false, false
		err = insertContent(moduleFile, func(line string) (string, int, bool) {
			line = strings.TrimSpace(line)
			if !moduleSetStart && strings.HasPrefix(line, "var Set = wire.NewSet") {
				moduleSetStart = true
			}
			if moduleSetStart && strings.HasSuffix(line, "// end") {
				moduleSetStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf(`wire.Struct(new(dao.%sRepo), "*"),\n`, cfg.Name))
				buf.WriteString(fmt.Sprintf(`wire.Struct(new(biz.%sBiz), "*"),\n`, cfg.Name))
				buf.WriteString(fmt.Sprintf(`wire.Struct(new(api.%sAPI), "*"),`, cfg.Name))
				return buf.String(), -1, true
			}

			if !moduleStructStart && strings.HasPrefix(line, fmt.Sprintf("type %s struct", cfg.ModuleName)) {
				moduleStructStart = true
			}
			if moduleStructStart && strings.HasSuffix(line, "// end") {
				moduleStructStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("%sAPI  *api.%sAPI", cfg.Name, cfg.Name))
				return buf.String(), -1, true
			}

			if !moduleMigrateStart && strings.HasPrefix(line, "return a.DB.AutoMigrate") {
				moduleMigrateStart = true
			}
			if moduleMigrateStart && strings.HasSuffix(line, "// end") {
				moduleMigrateStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("&typed.%s{},", cfg.Name))
				return buf.String(), -1, true
			}

			if !moduleGroupStart && strings.HasPrefix(line, `v1 := r.Group("v1")`) {
				moduleGroupStart = true
			}
			if moduleGroupStart && strings.HasSuffix(line, "// end") {
				moduleGroupStart = false
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf(`g%s := v1.Group("%s")\n`, cfg.Name, cfg.LowerPluralName))
				buf.WriteString("{\n")
				buf.WriteString(fmt.Sprintf(`g%s.GET("", a.%sAPI.Query)\n`, cfg.Name, cfg.Name))
				buf.WriteString(fmt.Sprintf(`g%s.GET(":id", a.%sAPI.Get)\n`, cfg.Name, cfg.Name))
				buf.WriteString(fmt.Sprintf(`g%s.POST("", a.%sAPI.Create)\n`, cfg.Name, cfg.Name))
				buf.WriteString(fmt.Sprintf(`g%s.PUT(":id", a.%sAPI.Update)\n`, cfg.Name, cfg.Name))
				buf.WriteString(fmt.Sprintf(`g%s.DELETE(":id", a.%sAPI.Delete)\n`, cfg.Name, cfg.Name))
				buf.WriteString("}")
				return buf.String(), -1, true
			}
			return "", 0, false
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func parseTplAndWriteFile(ctx context.Context, cfg *Config, result string) error {
	targetFile := filepath.Join(cfg.Directory, "internal", "module", cfg.ModuleLowerName, result, cfg.LowerUnderlineName+".go")
	if result == "module" {
		targetFile = filepath.Join(cfg.Directory, "internal", "module", cfg.ModuleLowerName, cfg.ModuleLowerName+".go")
	}

	if existsFile(targetFile) {
		logger.Context(ctx).Info("Skip file", zap.String("name", targetFile))
		return nil
	}

	logger.Context(ctx).Info("Generate file", zap.String("name", targetFile))
	tplFile := filepath.Join(cfg.TplDirectory, result+".go.tpl")
	tplBuf, err := readFile(tplFile)
	if err != nil {
		return err
	}

	tplResult, err := execParseTpl(tplBuf.String(), cfg)
	if err != nil {
		return err
	}

	err = writeFile(targetFile, tplResult)
	if err != nil {
		return err
	}

	if err := execGoFmt(targetFile); err != nil {
		logger.Context(ctx).Error("Failed to gofmt", zap.String("name", targetFile), zap.Error(err))
	}

	return nil
}
