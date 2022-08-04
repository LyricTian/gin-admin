package generate

import (
	"testing"

	"github.com/LyricTian/gin-admin/v9/pkg/util/yaml"
)

func TestConfig(t *testing.T) {
	cfg := &Config{
		PkgName:     "github.com/LyricTian/gin-admin/v9",
		ModelName:   "SYS",
		Name:        "Dictionary",
		Description: "Dictionary management for key/value pairs",
		Fields: []Field{
			{
				Name:                "ID",
				Type:                "string",
				Optional:            false,
				GormTag:             "size:20;primarykey;",
				BindingTag:          "",
				Comments:            "",
				InQuery:             false,
				QueryComments:       "",
				QueryIfExpression:   "",
				QueryWhereCondition: "",
				InCreate:            false,
			},
			{
				Name:                "Ns",
				Type:                "string",
				Optional:            false,
				GormTag:             "size:64;index;",
				BindingTag:          "required",
				Comments:            "Namespace of the dictionary",
				InQuery:             true,
				QueryComments:       "query namespace",
				QueryIfExpression:   `v := params.Namespace; len(v) > 0`,
				QueryWhereCondition: `"ns=?", v`,
				InCreate:            true,
			},
			{
				Name:                "Key",
				Type:                "string",
				Optional:            false,
				GormTag:             "size:64;index;",
				BindingTag:          "required",
				Comments:            "Key of the dictionary",
				InQuery:             true,
				QueryComments:       "query key",
				QueryIfExpression:   `v := params.Key; len(v) > 0`,
				QueryWhereCondition: `"key=?", v`,
				InCreate:            true,
			},
			{
				Name:                "Value",
				Type:                "string",
				Optional:            true,
				GormTag:             "size:4096;",
				BindingTag:          "",
				Comments:            "Value of the dictionary",
				InQuery:             false,
				QueryComments:       "",
				QueryIfExpression:   ``,
				QueryWhereCondition: ``,
				InCreate:            true,
			},
			{
				Name:                "CreatedAt",
				Type:                "time.Time",
				Optional:            false,
				GormTag:             "index;",
				BindingTag:          "",
				Comments:            "",
				InQuery:             false,
				QueryComments:       "",
				QueryIfExpression:   ``,
				QueryWhereCondition: ``,
				InCreate:            false,
			},
			{
				Name:                "CreatedBy",
				Type:                "string",
				Optional:            false,
				GormTag:             "size:20;",
				BindingTag:          "",
				Comments:            "",
				InQuery:             false,
				QueryComments:       "",
				QueryIfExpression:   ``,
				QueryWhereCondition: ``,
				InCreate:            false,
			},
			{
				Name:                "UpdatedAt",
				Type:                "time.Time",
				Optional:            false,
				GormTag:             "index;",
				BindingTag:          "",
				Comments:            "",
				InQuery:             false,
				QueryComments:       "",
				QueryIfExpression:   ``,
				QueryWhereCondition: ``,
				InCreate:            false,
			},
			{
				Name:                "UpdatedBy",
				Type:                "string",
				Optional:            false,
				GormTag:             "size:20;",
				BindingTag:          "",
				Comments:            "",
				InQuery:             false,
				QueryComments:       "",
				QueryIfExpression:   ``,
				QueryWhereCondition: ``,
				InCreate:            false,
			},
		},
	}

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + string(buf) + "\n")
}
