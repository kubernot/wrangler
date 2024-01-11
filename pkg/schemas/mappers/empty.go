package mappers

import (
	"github.com/kubernot/wrangler/pkg/data"
	"github.com/kubernot/wrangler/pkg/schemas"
)

type EmptyMapper struct {
}

func (e *EmptyMapper) FromInternal(data data.Object) {
}

func (e *EmptyMapper) ToInternal(data data.Object) error {
	return nil
}

func (e *EmptyMapper) ModifySchema(schema *schemas.Schema, schemas *schemas.Schemas) error {
	return nil
}
