package boilingcore

import (
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/importers"
)

// Config for the running of the commands
type Config struct {
	DriverName   string         `toml:"driver_name,omitempty" json:"driver_name,omitempty"`
	DriverConfig drivers.Config `toml:"driver_config,omitempty" json:"driver_config,omitempty"`

	PkgName          string   `toml:"pkg_name,omitempty" json:"pkg_name,omitempty"`
	OutFolder        string   `toml:"out_folder,omitempty" json:"out_folder,omitempty"`
	BaseDir          string   `toml:"base_dir,omitempty" json:"base_dir,omitempty"`
	Tags             []string `toml:"tags,omitempty" json:"tags,omitempty"`
	Replacements     []string `toml:"replacements,omitempty" json:"replacements,omitempty"`
	Debug            bool     `toml:"debug,omitempty" json:"debug,omitempty"`
	AddGlobal        bool     `toml:"add_global,omitempty" json:"add_global,omitempty"`
	AddPanic         bool     `toml:"add_panic,omitempty" json:"add_panic,omitempty"`
	NoContext        bool     `toml:"no_context,omitempty" json:"no_context,omitempty"`
	NoTests          bool     `toml:"no_tests,omitempty" json:"no_tests,omitempty"`
	NoHooks          bool     `toml:"no_hooks,omitempty" json:"no_hooks,omitempty"`
	NoAutoTimestamps bool     `toml:"no_auto_timestamps,omitempty" json:"no_auto_timestamps,omitempty"`
	NoRowsAffected   bool     `toml:"no_rows_affected,omitempty" json:"no_rows_affected,omitempty"`
	Wipe             bool     `toml:"wipe,omitempty" json:"wipe,omitempty"`
	StructTagCasing  string   `toml:"struct_tag_casing,omitempty" json:"struct_tag_casing,omitempty"`

	Imports importers.Collection `toml:"imports,omitempty" json:"imports,omitempty"`

	Aliases      Aliases       `toml:"aliases,omitempty" json:"aliases,omitempty"`
	TypeReplaces []TypeReplace `toml:"type_replaces,omitempty" json:"type_replaces,omitempty"`
}

// TypeReplace replaces a column type with something else
type TypeReplace struct {
	Match   drivers.Column `toml:"match,omitempty" json:"match,omitempty"`
	Replace drivers.Column `toml:"replace,omitempty" json:"replace,omitempty"`
	Imports importers.Set  `toml:"imports,omitempty" json:"imports,omitempty"`
}

// ConvertAliases is necessary because viper
func ConvertAliases(i interface{}) (a Aliases) {
	if i == nil {
		return a
	}

	topLevel := i.(map[string]interface{})
	tables := topLevel["tables"].(map[string]interface{})

	for name, tIntf := range tables {
		if a.Tables == nil {
			a.Tables = make(map[string]TableAlias)
		}

		t := tIntf.(map[string]interface{})

		var ta TableAlias

		if s := t["up_plural"]; s != nil {
			ta.UpPlural = s.(string)
		}
		if s := t["up_singular"]; s != nil {
			ta.UpSingular = s.(string)
		}
		if s := t["down_plural"]; s != nil {
			ta.DownPlural = s.(string)
		}
		if s := t["down_singular"]; s != nil {
			ta.DownSingular = s.(string)
		}

		if colsIntf, ok := t["columns"]; ok {
			cols := colsIntf.(map[string]interface{})
			ta.Columns = make(map[string]string)
			for k, v := range cols {
				ta.Columns[k] = v.(string)
			}
		}

		a.Tables[name] = ta
	}
	relationships := topLevel["relationships"].(map[string]interface{})

	for name, rIntf := range relationships {
		if a.Relationships == nil {
			a.Relationships = make(map[string]RelationshipAlias)
		}

		var ra RelationshipAlias
		rel := rIntf.(map[string]interface{})

		if s := rel["local"]; s != nil {
			ra.Local = s.(string)
		}
		if s := rel["foreign"]; s != nil {
			ra.Foreign = s.(string)
		}

		if len(ra.Foreign) == 0 || len(ra.Local) == 0 {
			panic("when defining a relationship alias, must name both sides of relationship")
		}

		a.Relationships[name] = ra
	}

	return a
}

// ConvertTypeReplace is necessary because viper
func ConvertTypeReplace(i interface{}) []TypeReplace {
	if i == nil {
		return nil
	}

	intfArray := i.([]interface{})
	var replaces []TypeReplace
	for _, r := range intfArray {
		replaceIntf := r.(map[string]interface{})
		replace := TypeReplace{}

		if replaceIntf["match"] == nil || replaceIntf["replace"] == nil {
			panic("replace types must specify both match and replace")
		}

		replace.Match = columnFromInterface(replaceIntf["match"])
		replace.Replace = columnFromInterface(replaceIntf["replace"])

		if imps := replaceIntf["imports"]; imps != nil {
			var err error
			replace.Imports, err = importers.SetFromInterface(imps)
			if err != nil {
				panic(err)
			}
		}

		replaces = append(replaces, replace)
	}

	return replaces
}

func columnFromInterface(i interface{}) (col drivers.Column) {
	m := i.(map[string]interface{})
	if s := m["name"]; s != nil {
		col.Name = s.(string)
	}
	if s := m["type"]; s != nil {
		col.Type = s.(string)
	}
	if s := m["db_type"]; s != nil {
		col.DBType = s.(string)
	}
	if s := m["udt_name"]; s != nil {
		col.UDTName = s.(string)
	}
	if s := m["full_db_type"]; s != nil {
		col.FullDBType = s.(string)
	}
	if s := m["arr_type"]; s != nil {
		col.ArrType = new(string)
		*col.ArrType = s.(string)
	}
	if b := m["auto_generated"]; b != nil {
		col.AutoGenerated = b.(bool)
	}
	if b := m["nullable"]; b != nil {
		col.Nullable = b.(bool)
	}

	return col
}
