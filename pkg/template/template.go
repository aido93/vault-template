package template

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/actano/vault-template/pkg/api"
	"os"
    "log"
	"strings"
	"text/template"
    "gopkg.in/yaml.v2"
)

type VaultTemplateRenderer struct {
	vaultClient api.VaultClient
}

func NewVaultTemplateRenderer(vaultToken, vaultEndpoint string) (*VaultTemplateRenderer, error) {
	vaultClient, err := api.NewVaultClient(vaultEndpoint, string(vaultToken))

	if err != nil {
		return nil, err
	}

	return &VaultTemplateRenderer{
		vaultClient: vaultClient,
	}, nil
}

func (v *VaultTemplateRenderer) RenderTemplate(templateContent string) (string, error) {
	funcMap := template.FuncMap{
		"vault":    v.vaultClient.QuerySecret,
		"vaultMap": v.vaultClient.QuerySecretMap,
        "toYaml":   toYaml,
        "fromYaml": fromYaml,
	}

	tmpl, err := template.
		New("template").
		Funcs(sprig.TxtFuncMap()).
		Funcs(funcMap).
		Parse(templateContent)

	if err != nil {
		return "", err
	}

	var outputBuffer bytes.Buffer

	envMap := envToMap()
	if err := tmpl.Execute(&outputBuffer, envMap); err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}

func envToMap() map[string]string {
	envMap := map[string]string{}

	for _, v := range os.Environ() {
		splitV := strings.Split(v, "=")
		envMap[splitV[0]] = splitV[1]
	}

	return envMap
}

func toYaml(src interface{}) string {
    data, err:= yaml.Marshal(src)
    if err!=nil {
        log.Print(err)
        return ""
    }
    return strings.TrimSuffix(string(data), "\n")
}

func fromYaml(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}
