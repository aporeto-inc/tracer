package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	type args struct {
		from  string
		to    string
		since time.Duration
	}
	tests := []struct {
		name              string
		args              args
		wantFromTime      time.Time
		wantToTime        time.Time
		wantSinceDuration time.Duration
		wantErr           bool
	}{
		{
			"Parse invalid from",
			args{
				from:  "chien",
				to:    "2020-10-22T17:56:17Z",
				since: time.Second,
			},
			time.Time{},
			time.Time{},
			0 * time.Second,
			true,
		},
		{
			"Parse invalid to",
			args{
				from:  "2020-10-22T17:56:17Z",
				to:    "chien",
				since: time.Second,
			},
			time.Time{},
			time.Time{},
			0 * time.Second,
			true,
		},
		{
			"Parse from and to",
			args{
				from:  "2020-10-22T17:56:16Z",
				to:    "2020-10-22T17:56:17Z",
				since: time.Second,
			},
			func() time.Time {
				t, err := time.Parse(time.RFC3339, "2020-10-22T17:56:16Z")
				if err != nil {
					panic(err)
				}
				return t
			}(),
			func() time.Time {
				t, err := time.Parse(time.RFC3339, "2020-10-22T17:56:17Z")
				if err != nil {
					panic(err)
				}
				return t
			}(),
			1 * time.Second,
			false,
		},
		{
			"Parse since and to",
			args{
				to:    "2020-10-22T17:56:17Z",
				since: time.Second,
			},
			func() time.Time {
				t, err := time.Parse(time.RFC3339, "2020-10-22T17:56:16Z")
				if err != nil {
					panic(err)
				}
				return t
			}(),
			func() time.Time {
				t, err := time.Parse(time.RFC3339, "2020-10-22T17:56:17Z")
				if err != nil {
					panic(err)
				}
				return t
			}(),
			1 * time.Second,
			false,
		},
		{
			"Parse since and from",
			args{
				from:  "2020-10-22T17:56:16Z",
				since: time.Second,
			},
			func() time.Time {
				t, err := time.Parse(time.RFC3339, "2020-10-22T17:56:16Z")
				if err != nil {
					panic(err)
				}
				return t
			}(),
			time.Now().Round(time.Second),
			func() time.Duration {
				t, err := time.Parse(time.RFC3339, "2020-10-22T17:56:16Z")
				if err != nil {
					panic(err)
				}
				return time.Now().Round(time.Second).Sub(t)
			}(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFromTime, gotToTime, gotSinceDuration, err := ParseTime(tt.args.from, tt.args.to, tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFromTime, tt.wantFromTime) {
				t.Errorf("ParseTime() gotFromTime = %v, want %v", gotFromTime, tt.wantFromTime)
			}
			if !reflect.DeepEqual(gotToTime, tt.wantToTime) {
				t.Errorf("ParseTime() gotToTime = %v, want %v", gotToTime, tt.wantToTime)
			}
			if gotSinceDuration != tt.wantSinceDuration {
				t.Errorf("ParseTime() gotSinceDuration = %v, want %v", gotSinceDuration, tt.wantSinceDuration)
			}
		})
	}
}
