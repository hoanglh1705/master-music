package converter

import "github.com/jinzhu/copier"

type (
	ModelConverter interface {
		FromModel(to interface{}, from interface{})
		ToModel(to interface{}, from interface{})
	}

	ObjectCopier interface {
		Copy(to interface{}, from interface{})
	}

	modelConverter struct {
	}

	adapterConverter struct {
	}

	objectCopier struct {
	}
)

func NewModelConverter() ModelConverter {
	return &modelConverter{}
}

func (h *modelConverter) FromModel(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}

func (h *modelConverter) ToModel(to interface{}, from interface{}) {
	_ = copier.Copy(to, from)
}
