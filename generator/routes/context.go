package routes

import (
	"github.com/haimkastner/gleece/cmd"
	"github.com/haimkastner/gleece/definitions"
)

type Argument struct {
	Type      definitions.ParamType
	Name      string
	ValueType any
}

type RouteCtx struct {
	definitions.RouteMetadata
}

type ControllerMeta struct {
	Routes []RouteCtx
}

type PackageImport struct {
	FullPath string
	Alias    string
}

type RoutesContext struct {
	PackageName string
	Controllers []definitions.ControllerMetadata
	//Imports     []PackageImport
	//Routes      []RouteCtx
}

func GetTemplateContext(
	config cmd.RoutesConfig,
	controllers []definitions.ControllerMetadata,
) (RoutesContext, error) {
	ctx := RoutesContext{Controllers: controllers}
	if len(config.PackageName) > 0 {
		ctx.PackageName = config.PackageName
	} else {
		ctx.PackageName = "routes"
	}

	/*
		imports := MapSet.NewSet[string]()

		for _, controller := range metadata.Controllers {
			// First, add the controller import
			imports.Add(controller.FullyQualifiedPackage)
		}
	*/
	return ctx, nil
}
