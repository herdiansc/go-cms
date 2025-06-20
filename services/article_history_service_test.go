package services

import (
	"errors"
	"net/url"
	"testing"

	"github.com/herdiansc/go-cms/models"
)

type mockArticleHistoryLister struct {
	d []models.ArticleHistory
	e error
}

func (m mockArticleHistoryLister) List(params map[string]interface{}) ([]models.ArticleHistory, error) {
	return m.d, m.e
}

var (
	mockSuccessArticleHistoryLister = mockArticleHistoryLister{
		d: []models.ArticleHistory{
			{
				Base:      models.Base{},
				Article:   "a",
				Version:   0,
				Status:    "",
				ArticleID: 0,
				Action:    "",
			},
		},
		e: nil,
	}
	mockEmptyArticleHistoryLister = mockArticleHistoryLister{
		d: []models.ArticleHistory{},
		e: nil,
	}
)

func TestListArticleHistoryServices_List(t *testing.T) {
	type fields struct {
		authData    any
		articleRepo mockArticleDetailer
		repo        mockArticleHistoryLister
	}
	type args struct {
		articleID int64
		q         url.Values
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
				articleRepo: mockSuccessArticleDetailer,
				repo:        mockSuccessArticleHistoryLister,
			},
			args: args{
				articleID: 1,
				q:         map[string][]string{"a": {"b"}},
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData:    "invalid",
				articleRepo: mockSuccessArticleDetailer,
				repo:        mockSuccessArticleHistoryLister,
			},
			args: args{
				articleID: 1,
				q:         map[string][]string{"a": {"b"}},
			},
			want: 400,
		},
		{
			name: "Failed to get data",
			fields: fields{
				authData:    mockValidAuthData,
				articleRepo: mockFailedArticleDetailer,
				repo:        mockSuccessArticleHistoryLister,
			},
			args: args{
				articleID: 1,
				q:         map[string][]string{"a": {"b"}},
			},
			want: 404,
		},
		{
			name: "Failed to get history data",
			fields: fields{
				authData:    mockValidAuthData,
				articleRepo: mockSuccessArticleDetailer,
				repo:        mockEmptyArticleHistoryLister,
			},
			args: args{
				articleID: 1,
				q:         map[string][]string{"a": {"b"}},
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewListArticleHistoryServices(
				tt.fields.authData,
				tt.fields.articleRepo,
				tt.fields.repo,
			)
			got, _ := svc.List(tt.args.articleID, tt.args.q)
			if got != tt.want {
				t.Errorf("ListArticleHistoryServices.List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockArticleHistoryDetailer struct {
	d models.ArticleHistory
	e error
}

func (m mockArticleHistoryDetailer) FindByParam(param string, value any) (models.ArticleHistory, error) {
	return m.d, m.e
}

var (
	mockSuccessArticleHistoryDetailer = mockArticleHistoryDetailer{
		d: models.ArticleHistory{},
		e: nil,
	}
	mockFailedArticleHistoryDetailer = mockArticleHistoryDetailer{
		d: models.ArticleHistory{},
		e: errors.New("error"),
	}
)

func TestDetailArticleHistoryServices_GetDetailByUUID(t *testing.T) {
	type fields struct {
		authData any
		repo     mockArticleHistoryDetailer
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
				repo:     mockSuccessArticleHistoryDetailer,
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
				repo:     mockSuccessArticleHistoryDetailer,
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
				repo:     mockFailedArticleHistoryDetailer,
			},
			args: args{
				id: 1,
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDetailArticleHistoryServices(tt.fields.authData, tt.fields.repo)
			got, _ := svc.GetDetailByUUID(tt.args.id)
			if got != tt.want {
				t.Errorf("DetailArticleHistoryServices.GetDetailByUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
