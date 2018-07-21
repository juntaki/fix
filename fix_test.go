package fix

import (
	"reflect"
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
		want    interface{}
		wantErr bool
	}{
		{
			name: "",
			args: args{
				target: &Test{
					Sub: TestSub{
						Value: "test",
					},
				},
			},
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
			raw, err := encode(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var got interface{}
			if err := decode(raw, &got); (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encode/decode = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultOutputPath(t *testing.T) {
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
			want: "testdata/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultOutputPath(tt.args.funcName); got != tt.want {
				t.Errorf("DefaultOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetOutputPathFunc(t *testing.T) {
	SetOutputPathFunc(func(funcName string) string {
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

func TestFix(t *testing.T) {
	SetOutputPathFunc(DefaultOutputPath)
	test := &Test{
		Sub: TestSub{
			Value: "test",
		},
	}

	err := Fix(test)
	if err != nil {
		t.Fatal(err)
	}

	test2 := &Test{
		Sub: TestSub{
			Value: "diff",
		},
	}
	err = Fix(test2)
	if err == nil {
		t.Fatal(err)
	}

	if err.Error() == `Diff: {*fix.Test}.Sub.Value:
        -: "diff"
        +: "test"` {
		t.Fatal(err.Error())
	}
}
