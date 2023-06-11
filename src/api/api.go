package api

import (
	"kalisto/src/models"
)

type Api struct {
}

func New() *Api {
	return &Api{}
}

func (a *Api) SpecFromProto(path string) (models.Spec, error) {
	return models.Spec{
		Services: []models.Service{
			{
				Name: "BookStore",
				Methods: []models.Method{
					{
						Name: "GetBook",
						Kind: models.CommunicationKindSimple,
						RequestMessage: models.Message{
							Name: "GetBookRequest",
							Fields: []models.Field{
								{
									Name: "id",
									Type: models.DataTypeString,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func (a *Api) SendGrpc(request models.Request) (models.Response, error) {
	return models.Response{
		Body: `{"name": "My super book"}`,
	}, nil
}
