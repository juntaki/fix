package fix

import (
	"testing"
)

type Test struct {
	Sub          TestSub
	PublicString string
	//privateString string
	PublicMap map[int]int
	//privateMap    map[int]int
}

type TestSub struct {
	Value string
}

func Test_encode(t *testing.T) {
	type args struct {
		target interface{}
	}

	sub := TestSub{
		Value: "test",
	}
	publicMap := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7}

	tests := []struct {
		name    string
		args    args
		codec   Codec
		wantErr bool
	}{
		{
			name: "json",
			args: args{
				target: &Test{
					Sub:       sub,
					PublicMap: publicMap,
				},
			},
			codec:   JSON,
			wantErr: false,
		},
		{
			name: "pp",
			args: args{
				target: &Test{
					Sub:       sub,
					PublicMap: publicMap,
				},
			},
			codec:   PP,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0 ; i<100; i++ {
				val1, err := tt.codec.Marshal(tt.args.target)
				if (err != nil) != tt.wantErr {
					t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				val2, err := tt.codec.Marshal(tt.args.target)
				if (err != nil) != tt.wantErr {
					t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				err = tt.codec.Compare(val1, val2)
				if (err != nil) != tt.wantErr {
					t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
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

func TestFixPP(t *testing.T) {
	SetOutputPathFunc(DefaultOutputPath)
	test := &Test{
		Sub: TestSub{
			Value: "test",
		},
	}

	err := PP.Fix(test)
	if err != nil {
		t.Fatal(err)
	}

	test2 := &Test{
		Sub: TestSub{
			Value: "diff",
		},
	}
	err = PP.Fix(test2)
	if err == nil {
		t.Fatal(err)
	}

	message := err.Error()
	err = PP.Fix(message, "message")
	if err != nil {
		t.Fatal(err)
	}
}
