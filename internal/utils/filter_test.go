package utils

import (
	"reflect"
	"sort"
	"testing"

	"github.com/aporeto-inc/tracer/internal/monitoring"
)

func TestFilter(t *testing.T) {
	type args struct {
		codes    string
		services []string
		urls     []string
		results  monitoring.APIErrors
	}
	tests := []struct {
		name    string
		args    args
		want    monitoring.APIErrors
		wantErr bool
	}{
		{
			"fail code parsing",
			args{
				codes: "chien",
			},
			nil,
			true,
		},
		{
			"fail code parsing bad range",
			args{
				codes: "foo-200",
			},
			nil,
			true,
		},
		{
			"fail code parsing bad range",
			args{
				codes: "200-chien",
			},
			nil,
			true,
		},
		{
			"fail code parsing invalid range",
			args{
				codes: "200-100",
			},
			nil,
			true,
		},
		{
			"fail code parsing invalid range",
			args{
				codes: "5-100-200",
			},
			nil,
			true,
		},
		{
			"fail code parsing parts",
			args{
				codes: "100,-2,-1",
			},
			nil,
			true,
		},
		{
			"proper parsing",
			args{
				codes:    "100,200-204,",
				services: []string{"foo"},
				urls:     []string{"/foo"},
			},
			monitoring.APIErrors{},
			false,
		},
		{
			"filtering no match",
			args{
				codes:    "100,200-204,",
				services: []string{"foo"},
				urls:     []string{"/foo"},
				results: monitoring.APIErrors{
					monitoring.APIError{
						Code:    300,
						Service: "zob",
						URL:     "/zob",
					},
					monitoring.APIError{
						Code:    500,
						Service: "foo",
						URL:     "/bar",
					},
					monitoring.APIError{
						Code:    500,
						Service: "/bar",
						URL:     "/foo",
					}},
			},
			monitoring.APIErrors{},
			false,
		},
		{
			"no filter",
			args{
				results: monitoring.APIErrors{
					monitoring.APIError{
						Count:   0,
						Code:    300,
						Service: "zob",
						URL:     "/zob",
					},
					monitoring.APIError{
						Count:   1,
						Code:    500,
						Service: "foo",
						URL:     "/bar",
					},
					monitoring.APIError{
						Count:   2,
						Code:    500,
						Service: "/bar",
						URL:     "/foo",
					}},
			},
			monitoring.APIErrors{
				monitoring.APIError{
					Count:   0,
					Code:    300,
					Service: "zob",
					URL:     "/zob",
				},
				monitoring.APIError{
					Count:   1,
					Code:    500,
					Service: "foo",
					URL:     "/bar",
				},
				monitoring.APIError{
					Count:   2,
					Code:    500,
					Service: "/bar",
					URL:     "/foo",
				}},
			false,
		},
		{
			"filtering works",
			args{
				codes:    "100,200-300,",
				services: []string{"foo"},
				urls:     []string{"/foo"},
				results: monitoring.APIErrors{
					monitoring.APIError{
						Count:   0,
						Code:    500,
						Service: "zob",
						URL:     "/zob",
					},
					monitoring.APIError{
						Count:   0,
						Code:    300,
						Service: "zob",
						URL:     "/zob",
					},
					monitoring.APIError{
						Count:   1,
						Code:    300,
						Service: "foo",
						URL:     "/bar",
					},
					monitoring.APIError{
						Count:   2,
						Code:    300,
						Service: "/bar",
						URL:     "/foo",
					},
					monitoring.APIError{
						Count:   2,
						Code:    500,
						Service: "/bar",
						URL:     "/foo",
					}},
			},
			monitoring.APIErrors{
				monitoring.APIError{
					Count:   1,
					Code:    300,
					Service: "foo",
					URL:     "/bar",
				},
				monitoring.APIError{
					Count:   2,
					Code:    300,
					Service: "/bar",
					URL:     "/foo",
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.args.codes, tt.args.services, tt.args.urls, tt.args.results)
			// we sort it so it's concistent
			sort.Sort(monitoring.ByCount(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
