package template

import (
	"context"

	"fmt"

	. "github.com/dave/jennifer/jen"
	"github.com/devimteam/microgen/generator/write_strategy"
)

const (
	EndpointsSetName = "EndpointsSet"
)

type endpointsTemplate struct {
	info *GenerationInfo
}

func NewEndpointsTemplate(info *GenerationInfo) Template {
	return &endpointsTemplate{
		info: info,
	}
}

func endpointStructName(str string) string {
	return str + "Endpoint"
}

// Renders endpoints file.
//
//		// This file was automatically generated by "microgen" utility.
//		// DO NOT EDIT.
//		package stringsvc
//
//		import (
//			context "context"
//			endpoint "github.com/go-kit/kit/endpoint"
//		)
//
//		type Endpoints struct {
//			CountEndpoint endpoint.Endpoint
//		}
//
func (t *endpointsTemplate) Render(ctx context.Context) write_strategy.Renderer {
	f := NewFile("transport")
	f.HeaderComment(t.info.FileHeader)

	f.Comment(fmt.Sprintf("%s implements %s API and used for transport purposes.", EndpointsSetName, t.info.Iface.Name))
	f.Type().Id(EndpointsSetName).StructFunc(func(g *Group) {
		for _, signature := range t.info.Iface.Methods {
			g.Id(endpointStructName(signature.Name)).Qual(PackagePathGoKitEndpoint, "Endpoint")
		}
	}).Line()

	return f
}

func (endpointsTemplate) DefaultPath() string {
	return filenameBuilder(PathTransport, "endpoints")
}

func (t *endpointsTemplate) Prepare(ctx context.Context) error {
	return nil
}

func (t *endpointsTemplate) ChooseStrategy(ctx context.Context) (write_strategy.Strategy, error) {
	return write_strategy.NewCreateFileStrategy(t.info.AbsOutputFilePath, t.DefaultPath()), nil
}
