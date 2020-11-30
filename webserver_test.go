package webserver

import (
	"reflect"
	"testing"
)

func TestStringExist(t *testing.T) {
	go StartServer()
	type args struct {
		req []string
	}
	tests := []struct {
		name    string
		args    args
		want    []bool
		wantErr bool
	}{
		{
			name: "t1",
			args: args{[]string{"abc", "abc", "efg"}},
			want: []bool{false, true, false},
			wantErr: true,
		},
		{
			name: "t2",
			args: args{[]string{"abc", "bcde", "efgg"}},
			want: []bool{true, false, false},
			wantErr: true,
		},
		{
			name: "t3",
			args: args{[]string{"abcd", "bcd", "efg","hijk"}},
			want: []bool{false, false, false, false},
			wantErr: true,
		},
		{
			name: "t4",
			args: args{[]string{"hijk", "bcdef", "efgg"}},
			want: []bool{true, false, false},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		//清空内部缓存数据
		if i % 2 == 0 {
			Store.Lock()
			Store.Data = make(map[string]int,0)
			Store.Unlock()
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringExist(tt.args.req)
			if (err != nil) == tt.wantErr {
				t.Errorf("StringExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringExist() = %v, want %v", got, tt.want)
			}
		})
	}
}