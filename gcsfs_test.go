// Copyright 2026 Jeremy Edwards
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ufs

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/googleapi"
)

var (
	fakeUpdatedTime = mustTime("2006-01-02T15:04:05Z")
)

func TestGetChunkSize(t *testing.T) {
	tests := []struct {
		size int
		want int
	}{
		{
			size: 0,
			want: minChunkSize,
		},
		{
			size: minChunkSize,
			want: minChunkSize * 2,
		},
		{
			size: minChunkSize * 2,
			want: minChunkSize * 4,
		},
		{
			size: minChunkSize * 100000,
			want: googleapi.DefaultUploadChunkSize,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("getChunkSize(%d) = %d", tc.size, tc.want), func(t *testing.T) {
			t.Parallel()
			got := getChunkSize(tc.size)
			if tc.want != got {
				t.Errorf("getChunkSize(%d) want: '%d' got: '%d'", tc.size, tc.want, got)
			}
		})
	}
}

func TestGCSJoin(t *testing.T) {
	tests := []struct {
		parts []string
		want  string
	}{
		{
			parts: []string{"gs://"},
			want:  "gs://",
		},
		{
			parts: []string{"gs://bucket\\"},
			want:  "gs://bucket",
		},
		{
			parts: []string{"gs:\\\\bucket\\a"},
			want:  "gs://bucket/a",
		},
		{
			parts: []string{"gs://first"},
			want:  "gs://first",
		},
		{
			parts: []string{"gs://first", "a", "b", "c"},
			want:  "gs://first/a/b/c",
		},
		{
			parts: []string{"gs://first/a/b", "c"},
			want:  "gs://first/a/b/c",
		},
		{
			parts: []string{"gs://first\\a\\b", "c"},
			want:  "gs://first/a/b/c",
		},
		{
			parts: []string{"gs://first\\a\\b\\c"},
			want:  "gs://first/a/b/c",
		},
		{
			parts: []string{"gs://first/a/..", "b/../c"},
			want:  "gs://first/c",
		},
		{
			parts: []string{"", "b/c"},
			want:  "b/c",
		},
		{
			parts: []string{"", "b/../c"},
			want:  "c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()
			got := gcsJoin(tc.parts...)
			if d := cmp.Diff(tc.want, got); d != "" {
				t.Errorf("got %s, want %s diff(-want,+got):\n %v", got, tc.want, d)
			}
		})
	}
}

func TestParseGCSPath(t *testing.T) {
	tests := []struct {
		name       string
		wantBucket string
		wantObject string
	}{
		{
			name:       "gs://first",
			wantBucket: "first",
		},
		{
			name:       "gs://first/a/b/c",
			wantBucket: "first",
			wantObject: "a/b/c",
		},
		{
			name:       "gs://first/a",
			wantBucket: "first",
			wantObject: "a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotBucket, gotObj, err := parseGCSPath(tc.name, "test")
			if err != nil {
				t.Error(err)
			}
			if d := cmp.Diff(tc.wantBucket, gotBucket); d != "" {
				t.Errorf("got %s, want %s diff(-want,+got):\n %v", gotBucket, tc.wantBucket, d)
			}
			if d := cmp.Diff(tc.wantObject, gotObj); d != "" {
				t.Errorf("got %s, want %s diff(-want,+got):\n %v", gotObj, tc.wantObject, d)
			}
		})
	}
}

func TestParseGCSPathErrors(t *testing.T) {
	tests := []struct {
		name            string
		wantErrContains string
	}{
		{
			name:            "file://a/c",
			wantErrContains: "does not contain the gs:// prefix",
		},
		{
			name:            "gs://",
			wantErrContains: "does not have a bucket name",
		},
		{
			name:            "gs:///object",
			wantErrContains: "does not have a bucket name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotBucket, gotObj, err := parseGCSPath(tc.name, "test")
			if err == nil {
				t.Errorf("want error that contains '%s'", tc.wantErrContains)
			} else if !strings.Contains(err.Error(), tc.wantErrContains) {
				t.Errorf("got: '%s' want contains: '%s'", err, tc.wantErrContains)
			}
			if gotBucket != "" {
				t.Errorf("bucket: want: '', got: '%s'", gotBucket)
			}
			if gotObj != "" {
				t.Errorf("object: want: '', got: '%s'", gotObj)
			}
		})
	}
}

func TestGCSFS(t *testing.T) {
	ctx := context.Background()
	client := createStorage(t)
	testFileSystem(t, func(name string) (FS, error) {
		return newGCSFSWithClient(ctx, client, name)
	}, "gs://first")
}

func createStorage(tb testing.TB) *storage.Client {
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "first",
				Name:       "a",
				Updated:    fakeUpdatedTime,
			},
			Content: []byte("content: first/a"),
		},
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "first",
				Name:       "b",
				Updated:    fakeUpdatedTime,
			},
			Content: []byte("content: first/b"),
		},
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "first",
				Name:       "dir/c",
				Updated:    fakeUpdatedTime,
			},
			Content: []byte("content: first/dir/c"),
		},
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "first",
				Name:       "dir/d",
				Updated:    fakeUpdatedTime,
			},
			Content: []byte("content: first/dir/d"),
		},
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "first",
				Name:       "dir/sub/e",
				Updated:    fakeUpdatedTime,
			},
			Content: []byte("content: first/dir/sub/e"),
		},
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "first",
				Name:       "dir/sub/f",
				Updated:    fakeUpdatedTime,
			},
			Content: []byte("content: first/dir/sub/f"),
		},
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: "first",
				Name:       "dir/subdir/g",
				Updated:    fakeUpdatedTime,
			},
			Content: []byte("content: first/dir/subdir/g"),
		},
	})

	tb.Cleanup(server.Stop)

	return server.Client()
}
