package template

import (
	"context"

	. "github.com/dave/jennifer/jen"
	"github.com/devimteam/microgen/generator/write_strategy"
	"github.com/devimteam/microgen/util"
	"github.com/vetcher/go-astra/types"
)

const (
	serviceErrorLoggingStructName = "errorLoggingMiddleware"
)

type errorLoggingTemplate struct {
	info *GenerationInfo
}

func NewErrorLoggingTemplate(info *GenerationInfo) Template {
	return &errorLoggingTemplate{
		info: info,
	}
}

func (t *errorLoggingTemplate) Render(ctx context.Context) write_strategy.Renderer {
	f := NewFile("service")
	f.ImportAlias(t.info.SourcePackageImport, serviceAlias)
	f.HeaderComment(t.info.FileHeader)

	f.Comment("ErrorLoggingMiddleware writes to logger any error, if it is not nil.").
		Line().Func().Id(util.ToUpperFirst(serviceErrorLoggingStructName)).Params(Id(loggerVarName).Qual(PackagePathGoKitLog, "Logger")).Params(Id(MiddlewareTypeName)).
		Block(t.newRecoverBody(t.info.Iface))

	f.Line()

	// Render type logger
	f.Type().Id(serviceErrorLoggingStructName).Struct(
		Id(loggerVarName).Qual(PackagePathGoKitLog, "Logger"),
		Id(nextVarName).Qual(t.info.SourcePackageImport, t.info.Iface.Name),
	)

	// Render functions
	for _, signature := range t.info.Iface.Methods {
		f.Line()
		f.Add(t.recoverFunc(ctx, signature)).Line()
	}

	return f
}

func (errorLoggingTemplate) DefaultPath() string {
	return filenameBuilder(PathService, "error_logging")
}

func (t *errorLoggingTemplate) Prepare(ctx context.Context) error {
	return nil
}

func (t *errorLoggingTemplate) ChooseStrategy(ctx context.Context) (write_strategy.Strategy, error) {
	return write_strategy.NewCreateFileStrategy(t.info.AbsOutputFilePath, t.DefaultPath()), nil
}

func (t *errorLoggingTemplate) newRecoverBody(i *types.Interface) *Statement {
	return Return(Func().Params(
		Id(nextVarName).Qual(t.info.SourcePackageImport, i.Name),
	).Params(
		Qual(t.info.SourcePackageImport, i.Name),
	).BlockFunc(func(g *Group) {
		g.Return(Op("&").Id(serviceErrorLoggingStructName).Values(
			Dict{
				Id(loggerVarName): Id(loggerVarName),
				Id(nextVarName):   Id(nextVarName),
			},
		))
	}))
}

func (t *errorLoggingTemplate) recoverFunc(ctx context.Context, signature *types.Function) *Statement {
	return methodDefinition(ctx, serviceErrorLoggingStructName, signature).
		BlockFunc(t.recoverFuncBody(signature))
}

func (t *errorLoggingTemplate) recoverFuncBody(signature *types.Function) func(g *Group) {
	return func(g *Group) {
		g.Defer().Func().Params().Block(
			If(Id(nameOfLastResultError(signature)).Op("!=").Nil()).Block(
				Id(util.LastUpperOrFirst(serviceErrorLoggingStructName)).Dot(loggerVarName).Dot("Log").Call(
					Lit("method"), Lit(signature.Name),
					Lit("message"), Id(nameOfLastResultError(signature)),
				),
			),
		).Call()

		g.Return().Id(util.LastUpperOrFirst(serviceErrorLoggingStructName)).Dot(nextVarName).Dot(signature.Name).Call(paramNames(signature.Args))
	}
}
