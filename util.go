package ecspresso

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go/private/protocol/json/jsonutil"
	"github.com/google/go-jsonnet"
)

func marshalJSON(s interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	b, err := jsonutil.BuildJSON(s)
	if err != nil {
		return nil, err
	}
	json.Indent(&buf, b, "", "  ")
	buf.WriteString("\n")
	return &buf, nil
}

func MarshalJSON(s interface{}) ([]byte, error) {
	b, err := marshalJSON(s)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), err
}

func MarshalJSONString(s interface{}) string {
	b, _ := marshalJSON(s)
	return b.String()
}

func isLongArnFormat(a string) (bool, error) {
	an, err := arn.Parse(a)
	if err != nil {
		return false, err
	}
	rs := strings.Split(an.Resource, "/")
	switch rs[0] {
	case "container-instance", "service", "task":
		return len(rs) >= 3, nil
	default:
		return false, nil
	}
}

func (d *App) readDefinitionFile(path string) ([]byte, error) {
	switch filepath.Ext(path) {
	case jsonnetExt:
		vm := jsonnet.MakeVM()
		for k, v := range d.ExtStr {
			vm.ExtVar(k, v)
		}
		for k, v := range d.ExtCode {
			vm.ExtCode(k, v)
		}
		jsonStr, err := vm.EvaluateFile(path)
		if err != nil {
			return nil, err
		}
		return d.loader.ReadWithEnvBytes([]byte(jsonStr))
	}
	return d.loader.ReadWithEnv(path)
}
