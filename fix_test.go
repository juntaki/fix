package fix

import (
	"testing"
)

type Test struct {
	Sub TestSub
}

type TestSub struct {
	Value string
}

func Test_encode(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		codec   Codec
		want    interface{}
		wantErr bool
	}{
		{
			name: "json",
			args: args{
				target: &Test{
					Sub: TestSub{
						Value: "test",
					},
				},
			},
			codec: JSON,
			want: &Test{
				Sub: TestSub{
					Value: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "gob",
			args: args{
				target: &Test{
					Sub: TestSub{
						Value: "test",
					},
				},
			},
			codec: Gob,
			want: &Test{
				Sub: TestSub{
					Value: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.codec.Marshal(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDefaultOutputPath(t *testing.T) {
	type args struct {
		funcName   string
		additional []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "base",
			args: args{
				funcName:   "test",
				additional: []string{"a", "b", "c"},
			},
			want: "testdata/test_a_b_c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultOutputPath(tt.args.funcName, tt.args.additional...); got != tt.want {
				t.Errorf("DefaultOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetOutputPathFunc(t *testing.T) {
	SetOutputPathFunc(func(funcName string, additional ...string) string {
		return funcName
	})
	type args struct {
		funcName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "base",
			args: args{
				funcName: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := outputPath(tt.args.funcName); got != tt.want {
				t.Errorf("DefaultOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixJSON(t *testing.T) {
	SetOutputPathFunc(DefaultOutputPath)
	test := &Test{
		Sub: TestSub{
			Value: "test",
		},
	}

	err := JSON.Fix(test)
	if err != nil {
		t.Fatal(err)
	}

	test2 := &Test{
		Sub: TestSub{
			Value: "diff",
		},
	}
	err = JSON.Fix(test2)
	if err == nil {
		t.Fatal(err)
	}

	message := err.Error()
	err = JSON.Fix(message, "message")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFixGob(t *testing.T) {
	SetOutputPathFunc(DefaultOutputPath)
	test := &Test{
		Sub: TestSub{
			Value: "test",
		},
	}

	err := Gob.Fix(test)
	if err != nil {
		t.Fatal(err)
	}

	test2 := &Test{
		Sub: TestSub{
			Value: "diff",
		},
	}
	err = Gob.Fix(test2)
	if err == nil {
		t.Fatal(err)
	}

	message := err.Error()
	err = Gob.Fix(message, "message")
	if err != nil {
		t.Fatal(err)
	}
}
