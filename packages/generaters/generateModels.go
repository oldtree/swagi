package generaters

import (
	"strings"
	"swagi/packages/swag"

	gen "github.com/dave/jennifer/jen"
	log "github.com/sirupsen/logrus"
)

func FilterDefinitionsName(params string) string {
	defineList := strings.Split(params, "/")
	length := len(defineList)
	return defineList[length-1]
}

func GeneratModels(sw *swag.Swagi) {
	f := gen.NewFile("models")
	f.ImportName("github.com/gin-gonic/gin", "gin")
	f.Func().Id("init").Params().Block(
		gen.Qual("fmt", "Println").Call(gen.Lit("models init")),
	)

	itemStatementma := make(map[string]*gen.Statement, len(sw.Description.Definitions))
	for namekey, item := range sw.Description.Definitions {
		name := FilterDefinitionsName(namekey)
		newItem, keyName := GenerateStruct(&item, name)
		if newItem == nil {
			continue
		}
		if item.Description == "persistent" && keyName != "" {
			methodMap := GenerateDatabaseOperationMap(&item, name, keyName)
			for key, value := range methodMap {
				itemStatementma[key] = value
			}
		}
		itemStatementma[name] = newItem
	}
	for _, statemet := range itemStatementma {
		if statemet == nil {
			continue
		}
		f.Add(statemet)
	}
	log.Infof("%#v", f)
	log.Info("end generate models")
}

func RebuildElementName(name string) string {
	return strings.ToUpper(string(name[0])) + string(name[1:])
}

func GenerateStruct(schema *swag.Schema, structName string) (*gen.Statement, string) {
	if schema.Type == "enum" {
		return nil, ""
	} else if schema.Type == "object" {
		var keyName string
		ProperList := make([]gen.Code, len(schema.Properties))
		for elemName, pro := range schema.Properties {
			var porper *gen.Statement
			switch pro.Type {
			case "integer":
				switch pro.Format {
				case "int":
					porper = gen.Id(RebuildElementName(elemName)).Int().Tag(map[string]string{"json": strings.ToLower(elemName)})
				case "int8":
					porper = gen.Id(RebuildElementName(elemName)).Int8().Tag(map[string]string{"json": strings.ToLower(elemName)})
				case "int16":
					porper = gen.Id(RebuildElementName(elemName)).Int16().Tag(map[string]string{"json": strings.ToLower(elemName)})
				case "int32":
					porper = gen.Id(RebuildElementName(elemName)).Int32().Tag(map[string]string{"json": strings.ToLower(elemName)})
				case "int64":
					porper = gen.Id(RebuildElementName(elemName)).Int64().Tag(map[string]string{"json": strings.ToLower(elemName)})
				}
			case "float":
				porper = gen.Id(RebuildElementName(elemName)).Float64().Tag(map[string]string{"json": strings.ToLower(elemName)})
			case "float32":
				porper = gen.Id(RebuildElementName(elemName)).Float32().Tag(map[string]string{"json": strings.ToLower(elemName)})
			case "float64":
				porper = gen.Id(RebuildElementName(elemName)).Float64().Tag(map[string]string{"json": strings.ToLower(elemName)})
			case "string":
				porper = gen.Id(RebuildElementName(elemName)).String().Tag(map[string]string{"json": strings.ToLower(elemName)})
			case "boolean":
				porper = gen.Id(RebuildElementName(elemName)).Bool().Tag(map[string]string{"json": strings.ToLower(elemName)})
			}
			ProperList = append(ProperList, porper)
			if pro.Description == "key" {
				keyName = elemName
			}
		}
		content := gen.Type().Id(RebuildElementName(structName)).Struct(
			ProperList...,
		)
		return content, keyName
	} else {
		log.Error("unknown object type")
		return nil, ""
	}
}

func GenerateDatabaseOperationList(schema *swag.Schema, structName string, keyName string) []gen.Code {
	methodList := make([]gen.Code, 4)
	methodList = append(methodList,
		GenerateDatabaseInsert(schema, structName, keyName),
		GenerateDatabaseGet(schema, structName, keyName),
		GenerateDatabaseUpdate(schema, structName, keyName),
		GenerateDatabaseDelete(schema, structName, keyName),
	)
	return methodList
}

func GenerateDatabaseOperationMap(schema *swag.Schema, structName string, keyName string) map[string]*gen.Statement {
	methodMap := make(map[string]*gen.Statement, 4)
	methodMap[keyName+"-get"] = GenerateDatabaseGet(schema, structName, keyName)
	methodMap[keyName+"-update"] = GenerateDatabaseUpdate(schema, structName, keyName)
	methodMap[keyName+"-post"] = GenerateDatabaseInsert(schema, structName, keyName)
	methodMap[keyName+"-delete"] = GenerateDatabaseDelete(schema, structName, keyName)
	return methodMap
}

func GenerateDatabaseInsert(schema *swag.Schema, structName string, keyName string) *gen.Statement {
	c := gen.Func().Params(
		gen.Id(string(structName[0])).Id("* " + RebuildElementName(structName)),
	).Id("Insert" + RebuildElementName(structName) + "by" + strings.ToUpper(keyName)).Params(
		gen.Id("id").Int(),
	).String().Block(
		gen.Return(gen.Id("b").Op("+").Id("c")),
	)
	return c
}

func GenerateDatabaseGet(schema *swag.Schema, structName string, keyName string) *gen.Statement {
	c := gen.Func().Params(
		gen.Id(string(structName[0])).Id("* " + RebuildElementName(structName)),
	).Id("Get" + RebuildElementName(structName) + "by" + strings.ToUpper(keyName)).Params(
		gen.Id("id").Int(),
	).String().Block(
		gen.Return(gen.Id("b").Op("+").Id("c")),
	)
	return c
}

func GenerateDatabaseUpdate(schema *swag.Schema, structName string, keyName string) *gen.Statement {
	c := gen.Func().Params(
		gen.Id(string(structName[0])).Id("* " + RebuildElementName(structName)),
	).Id("Update" + RebuildElementName(structName) + "by" + strings.ToUpper(keyName)).Params(
		gen.Id("id").Int(),
	).String().Block(
		gen.Return(gen.Id("b").Op("+").Id("c")),
	)
	return c
}

func GenerateDatabaseDelete(schema *swag.Schema, structName string, keyName string) *gen.Statement {
	c := gen.Func().Params(
		gen.Id(string(structName[0])).Id("* " + RebuildElementName(structName)),
	).Id("Delete" + RebuildElementName(structName) + "by" + strings.ToUpper(keyName)).Params(
		gen.Id("id").Int(),
	).String().Block(
		gen.Return(gen.Id("b").Op("+").Id("c")),
	)
	return c
}
