package tomladapter

import (
  "encoding/json"

  "github.com/caddyserver/caddy/v2/caddyconfig"
  "github.com/pelletier/go-toml"
)

func init() {
  caddyconfig.RegisterAdapter("toml", Adapter{})
}

// Adapter converts a TOML caddy config to json
type Adapter struct{}

func (a Adapter) Adapt(body []byte, _ map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
  tree, err := toml.LoadBytes(body)
  if err != nil {
    return nil, nil, err
  }
  b, err := json.Marshal(tree.ToMap())
  return b, nil, err
}
