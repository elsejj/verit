package version

import "testing"

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want Version
	}{
		{
			name: "basic semantic version",
			in:   "1.2.3",
			want: Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name: "with prerelease",
			in:   "0.1.0-alpha",
			want: Version{Major: 0, Minor: 1, Patch: 0, Prerelease: "alpha"},
		},
		{
			name: "with build metadata",
			in:   "2.10.5+build.11",
			want: Version{Major: 2, Minor: 10, Patch: 5, Build: "build.11"},
		},
		{
			name: "with prerelease and build metadata",
			in:   "3.4.5-rc.1+sha.abcdef",
			want: Version{Major: 3, Minor: 4, Patch: 5, Prerelease: "rc.1", Build: "sha.abcdef"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v, want nil", tt.in, err)
			}

			if !got.Equal(&tt.want) {
				t.Fatalf("Parse(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseInvalid(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"1.0",
		"1.0.0.0",
		"01.2.3",
		"a.b.c",
	}

	for _, in := range tests {
		in := in
		t.Run("invalid_"+in, func(t *testing.T) {
			t.Parallel()

			if _, err := Parse(in); err == nil {
				t.Fatalf("Parse(%q) error = nil, want non-nil", in)
			}
		})
	}
}

func TestVersionComparison(t *testing.T) {
	t.Parallel()

	type cmpCase struct {
		name string
		a    string
		b    string
		want int // -1 if a<b, 0 if a==b, 1 if a>b
	}

	tests := []cmpCase{
		{
			name: "identical versions",
			a:    "1.2.3",
			b:    "1.2.3",
			want: 0,
		},
		{
			name: "patch greater",
			a:    "1.2.4",
			b:    "1.2.3",
			want: 1,
		},
		{
			name: "minor less",
			a:    "1.2.3",
			b:    "1.3.0",
			want: -1,
		},
		{
			name: "major greater",
			a:    "2.0.0",
			b:    "1.9.9",
			want: 1,
		},
		{
			name: "prerelease lexical comparison",
			a:    "1.2.3-alpha",
			b:    "1.2.3-beta",
			want: -1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, err := Parse(tt.a)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v, want nil", tt.a, err)
			}
			b, err := Parse(tt.b)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v, want nil", tt.b, err)
			}

			switch tt.want {
			case 1:
				if !a.GreaterThan(b) {
					t.Fatalf("%q should be greater than %q", tt.a, tt.b)
				}
				if a.LessThan(b) {
					t.Fatalf("%q should not be less than %q", tt.a, tt.b)
				}
				if a.Equal(b) {
					t.Fatalf("%q should not equal %q", tt.a, tt.b)
				}
			case -1:
				if !a.LessThan(b) {
					t.Fatalf("%q should be less than %q", tt.a, tt.b)
				}
				if a.GreaterThan(b) {
					t.Fatalf("%q should not be greater than %q", tt.a, tt.b)
				}
				if a.Equal(b) {
					t.Fatalf("%q should not equal %q", tt.a, tt.b)
				}
			case 0:
				if a.GreaterThan(b) || b.GreaterThan(a) {
					t.Fatalf("%q and %q should be equal", tt.a, tt.b)
				}
				if a.LessThan(b) || b.LessThan(a) {
					t.Fatalf("%q and %q should not be less than each other", tt.a, tt.b)
				}
				if !a.Equal(b) || !b.Equal(a) {
					t.Fatalf("%q and %q should be equal", tt.a, tt.b)
				}
			default:
				t.Fatalf("unsupported want value %d", tt.want)
			}
		})
	}
}
