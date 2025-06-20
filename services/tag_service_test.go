package services

import (
	"errors"
	"net/url"
	"testing"

	"github.com/herdiansc/go-cms/models"
)

type mockTagCreator struct {
	d models.Tag
	e error
}

func (m mockTagCreator) Create(data models.Tag) (models.Tag, error) {
	return m.d, m.e
}

var (
	mockValidAuthData = models.VerifyData{
		ID:       1,
		UUID:     "abc-123",
		Username: "test",
	}
	mockSuccessTagCreator = mockTagCreator{
		d: models.Tag{},
		e: nil,
	}
	mockFailedTagCreator = mockTagCreator{
		d: models.Tag{},
		e: errors.New("error"),
	}
)

func TestCreateTagServices_Create(t *testing.T) {
	type fields struct {
		authData  any
		decoder   mockJsonDecoder
		validator mockRequestValidator
		repo      mockTagCreator
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Positive",
			fields: fields{
				authData:  mockValidAuthData,
				decoder:   mockSuccessJsonDecoder,
				validator: mockSuccessRequestValidator,
				repo:      mockSuccessTagCreator,
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData:  "invalid",
				decoder:   mockSuccessJsonDecoder,
				validator: mockSuccessRequestValidator,
				repo:      mockSuccessTagCreator,
			},
			want: 400,
		},
		{
			name: "Failed to decode json data",
			fields: fields{
				authData:  mockValidAuthData,
				decoder:   mockFailedJsonDecoder,
				validator: mockSuccessRequestValidator,
				repo:      mockSuccessTagCreator,
			},
			want: 400,
		},
		{
			name: "Failed to validate data",
			fields: fields{
				authData:  mockValidAuthData,
				decoder:   mockSuccessJsonDecoder,
				validator: mockFailedRequestValidator,
				repo:      mockSuccessTagCreator,
			},
			want: 400,
		},
		{
			name: "Failed to save data",
			fields: fields{
				authData:  mockValidAuthData,
				decoder:   mockSuccessJsonDecoder,
				validator: mockSuccessRequestValidator,
				repo:      mockFailedTagCreator,
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCreateTagServices(tt.fields.authData, tt.fields.decoder, tt.fields.validator, tt.fields.repo)
			got, _ := svc.Create()
			if got != tt.want {
				t.Errorf("CreateTagServices.Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockTagLister struct {
	d []models.TagListItem
	e error
}

func (m mockTagLister) List(params map[string]interface{}) ([]models.TagListItem, error) {
	return m.d, m.e
}

var (
	mockSuccessTagLister = mockTagLister{
		d: []models.TagListItem{
			{
				ID:         1,
				Title:      "b",
				UsageCount: 1,
			},
		},
		e: nil,
	}
	mockEmptyTagLister = mockTagLister{
		d: []models.TagListItem{},
		e: nil,
	}
)

func TestListTagServices_List(t *testing.T) {
	type fields struct {
		authData any
		repo     mockTagLister
	}
	type args struct {
		q url.Values
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Positive",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockSuccessTagLister,
			},
			args: args{
				q: map[string][]string{"a": {"b"}},
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData: "invalid",
				repo:     mockSuccessTagLister,
			},
			args: args{},
			want: 400,
		},
		{
			name: "Positive",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockSuccessTagLister,
			},
			args: args{},
			want: 200,
		},
		{
			name: "Positive",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockEmptyTagLister,
			},
			args: args{
				q: map[string][]string{"a": {"b"}},
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewListTagServices(tt.fields.authData, tt.fields.repo)
			got, _ := svc.List(tt.args.q)
			if got != tt.want {
				t.Errorf("ListTagServices.List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockTagDetailer struct {
	d models.TagDetail
	e error
}

func (m mockTagDetailer) FindByParam(param string, value any) (models.TagDetail, error) {
	return m.d, m.e
}

var (
	mockSuccessTagDetailer = mockTagDetailer{
		d: models.TagDetail{},
		e: nil,
	}
	mockFailedTagDetailer = mockTagDetailer{
		d: models.TagDetail{},
		e: errors.New("error"),
	}
)

func TestDetailTagServices_GetDetailByUUID(t *testing.T) {
	type fields struct {
		authData any
		repo     mockTagDetailer
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Positive",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockSuccessTagDetailer,
			},
			args: args{},
			want: 200,
		},
		{
			name: "Failed to get data",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockFailedTagDetailer,
			},
			args: args{},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDetailTagServices(tt.fields.authData, tt.fields.repo)
			got, _ := svc.GetDetailByUUID(tt.args.id)
			if got != tt.want {
				t.Errorf("DetailTagServices.GetDetailByUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
