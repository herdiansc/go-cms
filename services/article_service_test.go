package services

import (
	"errors"
	"net/url"
	"testing"

	"github.com/herdiansc/go-cms/models"
)

type mockArticleProcessor struct {
	d models.Article
	e error
}

func (m mockArticleProcessor) Create(writerID int64, data models.CreateArticleRequest) (models.Article, error) {
	return m.d, m.e
}

type mockArticleHistoryCreator struct {
	e error
}

func (m mockArticleHistoryCreator) Create(action string, data models.Article) error {
	return m.e
}

var (
	mockSuccessArticleProcessor = mockArticleProcessor{
		d: models.Article{},
		e: nil,
	}
	mockFailedArticleProcessor = mockArticleProcessor{
		d: models.Article{},
		e: errors.New("error"),
	}
	mockSuccessArticleHistoryCreator = mockArticleHistoryCreator{
		e: nil,
	}
)

func TestCreateArticleServices_Create(t *testing.T) {
	type fields struct {
		authData    any
		decoder     mockJsonDecoder
		validator   mockRequestValidator
		repo        mockArticleProcessor
		historyRepo mockArticleHistoryCreator
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Positive",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockSuccessJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockSuccessArticleProcessor,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData:    "invalid",
				decoder:     mockSuccessJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockSuccessArticleProcessor,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			want: 400,
		},
		{
			name: "Failed to decode json data",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockFailedJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockSuccessArticleProcessor,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			want: 400,
		},
		{
			name: "Failed to validate data",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockSuccessJsonDecoder,
				validator:   mockFailedRequestValidator,
				repo:        mockSuccessArticleProcessor,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			want: 400,
		},
		{
			name: "Failed to save data",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockSuccessJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockFailedArticleProcessor,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			want: 500,
		},
		{
			name: "Positive",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockSuccessJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockSuccessArticleProcessor,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCreateArticleServices(
				tt.fields.authData,
				tt.fields.decoder,
				tt.fields.validator,
				tt.fields.repo,
				tt.fields.historyRepo,
			)
			got, _ := svc.Create()
			if got != tt.want {
				t.Errorf("CreateArticleServices.Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockArticleLister struct {
	d []models.Article
	e error
}

func (m mockArticleLister) List(params map[string]interface{}) ([]models.Article, error) {
	return m.d, m.e
}

var (
	mockSuccessArticleLister = mockArticleLister{
		d: []models.Article{
			{
				Base:     models.Base{},
				Title:    "a",
				Content:  "",
				Status:   "",
				WriterID: 0,
				Slug:     "",
			},
		},
		e: nil,
	}
	mockEmptyArticleLister = mockArticleLister{
		d: []models.Article{},
		e: nil,
	}
)

func TestListArticleServices_List(t *testing.T) {
	type fields struct {
		authData any
		repo     mockArticleLister
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
				repo:     mockSuccessArticleLister,
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
				repo:     mockSuccessArticleLister,
			},
			args: args{
				q: map[string][]string{"a": {"b"}},
			},
			want: 400,
		},
		{
			name: "Failed to get data",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockEmptyArticleLister,
			},
			args: args{
				q: map[string][]string{"a": {"b"}},
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewListArticleServices(tt.fields.authData, tt.fields.repo)
			got, _ := svc.List(tt.args.q)
			if got != tt.want {
				t.Errorf("ListArticleServices.List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockArticleDetailer struct {
	d models.Article
	e error
}

func (m mockArticleDetailer) FindByParam(param string, value any) (models.Article, error) {
	return m.d, m.e
}

var (
	mockSuccessArticleDetailer = mockArticleDetailer{
		d: models.Article{},
		e: nil,
	}
	mockFailedArticleDetailer = mockArticleDetailer{
		d: models.Article{},
		e: errors.New("error"),
	}
)

func TestDetailArticleServices_GetDetailByUUID(t *testing.T) {
	type fields struct {
		authData any
		repo     mockArticleDetailer
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
				repo:     mockSuccessArticleDetailer,
			},
			args: args{
				id: 1,
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData: "invalid",
				repo:     mockSuccessArticleDetailer,
			},
			args: args{
				id: 1,
			},
			want: 400,
		},
		{
			name: "Failed to get data",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockFailedArticleDetailer,
			},
			args: args{
				id: 1,
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDetailArticleServices(tt.fields.authData, tt.fields.repo)
			got, _ := svc.GetDetailByUUID(tt.args.id)
			if got != tt.want {
				t.Errorf("DetailArticleServices.GetDetailByUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockArticleDeleter struct {
	e error
}

func (m mockArticleDeleter) DeleteByParam(param string, value any) error {
	return m.e
}

var (
	mockSuccessArticleDeleter = mockArticleDeleter{
		e: nil,
	}
	mockFailedArticleDeleter = mockArticleDeleter{
		e: errors.New("error"),
	}
)

func TestDeleteArticleServices_Delete(t *testing.T) {
	type fields struct {
		authData any
		repo     mockArticleDeleter
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
				repo:     mockSuccessArticleDeleter,
			},
			args: args{
				id: 1,
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData: "invalid",
				repo:     mockSuccessArticleDeleter,
			},
			args: args{
				id: 1,
			},
			want: 400,
		},
		{
			name: "Failed to delete data",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockFailedArticleDeleter,
			},
			args: args{
				id: 1,
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDeleteArticleServices(tt.fields.authData, tt.fields.repo)
			got, _ := svc.Delete(tt.args.id)
			if got != tt.want {
				t.Errorf("DeleteArticleServices.Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockArticlePatcher struct {
	d models.Article
	e error
}

func (m mockArticlePatcher) PatchByParam(id int64, param string, value any) (models.Article, error) {
	return m.d, m.e
}

var (
	mockSuccessArticlePatcher = mockArticlePatcher{
		d: models.Article{},
		e: nil,
	}
	mockFailedArticlePatcher = mockArticlePatcher{
		d: models.Article{},
		e: errors.New("error"),
	}
)

func TestPatchArticleServices_Patch(t *testing.T) {
	type fields struct {
		authData    any
		decoder     mockJsonDecoder
		validator   mockRequestValidator
		repo        mockArticlePatcher
		historyRepo mockArticleHistoryCreator
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
				authData:    mockValidAuthData,
				decoder:     mockSuccessJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockSuccessArticlePatcher,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			args: args{
				id: 1,
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData:    "invalid",
				decoder:     mockSuccessJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockSuccessArticlePatcher,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			args: args{
				id: 1,
			},
			want: 400,
		},
		{
			name: "Failed to decode json data",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockFailedJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockSuccessArticlePatcher,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			args: args{
				id: 1,
			},
			want: 400,
		},
		{
			name: "Failed to validate data",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockSuccessJsonDecoder,
				validator:   mockFailedRequestValidator,
				repo:        mockSuccessArticlePatcher,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			args: args{
				id: 1,
			},
			want: 400,
		},
		{
			name: "Failed to save data",
			fields: fields{
				authData:    mockValidAuthData,
				decoder:     mockSuccessJsonDecoder,
				validator:   mockSuccessRequestValidator,
				repo:        mockFailedArticlePatcher,
				historyRepo: mockSuccessArticleHistoryCreator,
			},
			args: args{
				id: 1,
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPatchArticleServices(
				tt.fields.authData,
				tt.fields.decoder,
				tt.fields.validator,
				tt.fields.repo,
				tt.fields.historyRepo,
			)
			got, _ := svc.Patch(tt.args.id)
			if got != tt.want {
				t.Errorf("PatchArticleServices.Patch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
