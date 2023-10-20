package definitions

import (
	"fmt"
	"kalisto/src/models"
	"strings"
)

func MakeRequestExample(m models.Message, links map[string]models.Message) string {
	return makeRequestExample(make(map[string]bool), m, links, 2, "")
}

func makeRequestExample(set map[string]bool, m models.Message, links map[string]models.Message, space int, parent string) string {
	var buf strings.Builder
	buf.WriteString("{\n")

	for _, field := range m.Fields {
		setKey := fmt.Sprintf("%s:%s:%s", m.FullName, parent, field.FullName)
		if field.Message != "" && set[setKey] {
			continue
		}
		set[setKey] = true
		v := makeExampleValue(set, field, links, space, m.FullName)
		if v == "" {
			continue
		}

		tpl := "%s%s: %s,\n"
		if field.Type == models.DataTypeOneOf {
			tpl = "%s%s: %s"
		}
		line := fmt.Sprintf(tpl, strings.Repeat(" ", space), field.Name, v)
		buf.WriteString(line)
	}

	closeBracket := fmt.Sprintf("%s}", strings.Repeat(" ", space-2))
	buf.WriteString(closeBracket)
	return buf.String()
}

func makeExampleValue(set map[string]bool, field models.Field, links map[string]models.Message, space int, parent string) string {
	if field.Repeated {
		fieldCp := field
		fieldCp.Repeated = false
		v := makeExampleValue(set, fieldCp, links, space, "")
		return fmt.Sprintf("[%s]", v)
	}

	var v string
	switch field.Type {
	case models.DataTypeString:
		v = `"string"`
	case models.DataTypeBool:
		v = `true`
	case models.DataTypeInt32, models.DataTypeInt64, models.DataTypeUint32, models.DataTypeUint64:
		v = `1`
	case models.DataTypeFloat32, models.DataTypeFloat64:
		v = `3.14`
	case models.DataTypeBytes:
		v = `"{json: true}"`
	case models.DataTypeEnum:
		if len(field.Enum) == 0 {
			return ""
		}
		v = fmt.Sprintf(`%d`, field.Enum[0])
	case models.DataTypeDuration:
		v = "1576800000000000"
	case models.DataTypeDate:
		v = "Date.now()"
	case models.DataTypeStruct:
		if field.MapKey != nil && field.MapValue != nil {
			key := makeExampleValue(set, *field.MapKey, links, space, "")
			value := makeExampleValue(set, *field.MapValue, links, space, "")
			v = fmt.Sprintf(`{%s: %s}`, key, value)
		} else {
			link, ok := links[field.Message]
			if !ok {
				return ""
			}
			if strings.HasPrefix(parent, field.FullName) {
				return ""
			}
			parent += ":" + field.FullName
			v = makeRequestExample(set, link, links, space+2, parent)
		}
	case models.DataTypeOneOf:
		var oneOfBuf strings.Builder
		for i, one := range field.OneOf {
			oneV := makeExampleValue(set, one, links, space, field.FullName)
			oneV = fmt.Sprintf("{%s: %s},\n", one.Name, oneV)
			if i != 0 {
				oneV = field.Name + ": " + oneV
				lines := strings.Split(oneV, "\n")
				for i, line := range lines {
					if strings.TrimSpace(line) == "" {
						continue
					}
					lines[i] = strings.Repeat(" ", space) + "// " + line
				}
				oneV = strings.Join(lines, "\n")
			}
			oneOfBuf.WriteString(oneV)
		}
		v = oneOfBuf.String()
	}

	return v
}
