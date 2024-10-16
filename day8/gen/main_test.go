package gen

import "testing"

func Test_gen(t *testing.T) {
	type args struct {
		w       *io.Writer
		srcFile string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen(tt.args.w, tt.args.srcFile)
		})
	}
}
