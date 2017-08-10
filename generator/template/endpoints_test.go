package template

import (
	goparser "go/parser"
	"go/token"
	"testing"

	"bytes"

	"github.com/devimteam/microgen/generator"
	parser "github.com/devimteam/microgen/parser"
)

func TestEndpointsForCountSvc(t *testing.T) {
	src := `package stringsvc

import (
	"context"
)

type StringService interface {
	Count(ctx context.Context, text string, symbol string) (count int, positions []int)
}`

	out := `// This file was automatically generated by "microgen" utility.
// Please, do not edit.
package stringsvc

import (
	context "context"
	endpoint "github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CountEndpoint endpoint.Endpoint
}

func (e *Endpoints) Count(ctx context.Context, text string, symbol string) (count int, positions []int) {
	req := CountRequest{
		Symbol: symbol,
		Text:   text,
	}
	resp, err := e.CountEndpoint(ctx, &req)
	if err != nil {
		return
	}
	return resp.(*CountResponse).Count, resp.(*CountResponse).Positions
}

func CountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*CountRequest)
		count, positions := svc.Count(ctx, req.Text, req.Symbol)
		return &CountResponse{
			Count:     count,
			Positions: positions,
		}, nil
	}
}
` // Blank line!
	f, err := goparser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		t.Errorf("unable to parse file: %v", err)
	}
	fs, err := parser.ParseInterface(f, "StringService")
	if err != nil {
		t.Errorf("could not get interface func signatures: %v", err)
	}
	buf := bytes.NewBuffer([]byte{})
	gen := generator.NewGenerator([]generator.Template{
		&EndpointsTemplate{},
	}, fs, generator.NewWriterStrategy(buf))
	err = gen.Generate()
	if err != nil {
		t.Errorf("unable to generate: %v", err)
	}
	if buf.String() != out {
		t.Errorf("Got:\n\n%v\n\nShould be:\n\n%v", buf.String(), out)
	}
}