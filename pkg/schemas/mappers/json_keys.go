package mappers

import (
	"github.com/kubernot/wrangler/pkg/data"
	"github.com/kubernot/wrangler/pkg/data/convert"
	types "github.com/kubernot/wrangler/pkg/schemas"
)

type JSONKeys struct {
}

func (d JSONKeys) FromInternal(data data.Object) {
}

func (d JSONKeys) ToInternal(data data.Object) error {
	for key, value := range data {
		newKey := convert.ToJSONKey(key)
		if newKey != key {
			data[newKey] = value
			delete(data, key)
		}
	}
	return nil
}

func (d JSONKeys) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
