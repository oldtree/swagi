package generaters

import (
	"strings"
	"swagi/packages/swag"

	gen "github.com/dave/jennifer/jen"
	_ "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var requestApplicationType string
var responseApplicationType string

func GeneratRouters(sw *swag.Swagi) {
	log.Info("start route generate")
	f := gen.NewFile("handles")
	f.ImportName("github.com/gin-gonic/gin", "gin")
	f.ImportName("net/http", "http")
	f.Func().Id("init").Params().Block(
		gen.Qual("fmt", "Println").Call(gen.Lit("router init")),
	)
	var FuncNameMap = make(map[string]map[string]string)
	BuildHandlers(sw, f, FuncNameMap)
	f.Func().Id("Handles").Params(gen.List(gen.Id("app")).Op(" *").Id(gen.Qual("gin", "Engine").GoString())).Block(
		gen.Qual("log", "Info").Call(gen.Lit("init handles")),
		gen.List(gen.Id("handleGroup")).Op(":=").Id("app").Dot(`Group("/")`),
		InstallRouters(f),
	)

	log.Infof("%#v", f)
	log.Info("route generate end")
}

func InstallRouters(stmt *gen.File) *gen.Statement {
	return nil
}

func BuildHandlers(sw *swag.Swagi, stmt *gen.File, funcMap map[string]map[string]string) {
	routePath := sw.Description.Paths
	var routeGroup = make(map[string]map[string]*gen.Statement, len(routePath))
	for pathName, pathItem := range routePath {
		routeGroup[pathName] = make(map[string]*gen.Statement)
		if pathItem.Delete != nil {

			//routeGroup[pathName]["delete"] = BuildDeleteMethod(pathItem.Delete, pathName, funcMap)
		}
		if pathItem.Get != nil {
			//routeGroup[pathName]["get"] = BuildGetMethod(pathItem.Get, pathName, funcMap)
		}
		if pathItem.Put != nil {
			//routeGroup[pathName]["put"] = BuildPutMethod(pathItem.Put, pathName, funcMap)
		}
		if pathItem.Post != nil {
			//routeGroup[pathName]["post"] = BuildPostMethod(pathItem.Post, pathName, funcMap)
		}
	}
	for _, methodList := range routeGroup {
		for _, handleStmt := range methodList {
			stmt.Add(handleStmt)
		}
	}
	return
}

func BuildFuncName(path string) string {
	return strings.ToUpper(string(path[0])) + string(path[1:])
}

func BuildPostMethod(item *swag.Operation, path string, funcMap map[string]string) *gen.Statement {
	var content gen.Statement
	content.Func().Id("Handle" + BuildFuncName(item.OperationID)).Params(
		gen.List(
			gen.Id("ctx").Op(" *").Id(gen.Qual("gin", "context").GoString()),
		),
	).Block()
	return &content
}

func BuildGetMethod(item *swag.Operation, path string, funcMap map[string]string) *gen.Statement {
	var content gen.Statement
	content.Func().Id("Handle" + BuildFuncName(item.OperationID)).Params(
		gen.List(
			gen.Id("ctx").Op(" *").Id(gen.Qual("gin", "context").GoString())),
	).Block()
	return &content
}

func BuildDeleteMethod(item *swag.Operation, path string, funcMap map[string]string) *gen.Statement {
	var content gen.Statement
	content.Func().Id("Handle" + BuildFuncName(item.OperationID)).Params(
		gen.List(
			gen.Id("ctx").Op(" *").Id(gen.Qual("gin", "context").GoString())),
	).Block()
	return &content
}

func BuildPutMethod(item *swag.Operation, path string, funcMap map[string]string) *gen.Statement {
	var content gen.Statement
	content.Func().Id("Handle" + BuildFuncName(item.OperationID)).Params(
		gen.List(
			gen.Id("ctx").Op(" *").Id(gen.Qual("gin", "context").GoString())),
	).Block()
	return &content
}
