package versioning

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in        string
		wantMajor int
		wantMinor int
		wantPatch int
		wantPre   string
		wantErr   bool
	}{
		{in: "1.2.3", wantMajor: 1, wantMinor: 2, wantPatch: 3},
		{in: "v0.1.0", wantMajor: 0, wantMinor: 1, wantPatch: 0},
		{in: "1.2.3-rc1", wantMajor: 1, wantMinor: 2, wantPatch: 3, wantPre: "rc1"},
		{in: "1.2.3+build.42", wantMajor: 1, wantMinor: 2, wantPatch: 3},
		{in: "1.2", wantErr: true},
		{in: "v1.x.0", wantErr: true},
		{in: "", wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got, err := Parse(c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, c.wantErr)
			}
			if c.wantErr {
				return
			}
			if got.Major != c.wantMajor || got.Minor != c.wantMinor || got.Patch != c.wantPatch {
				t.Errorf("triple = %d.%d.%d, want %d.%d.%d",
					got.Major, got.Minor, got.Patch,
					c.wantMajor, c.wantMinor, c.wantPatch)
			}
			if got.Prerelease != c.wantPre {
				t.Errorf("pre = %q, want %q", got.Prerelease, c.wantPre)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{a: "1.0.0", b: "1.0.0", want: 0},
		{a: "1.0.0", b: "0.9.9", want: 1},
		{a: "0.9.9", b: "1.0.0", want: -1},
		{a: "1.2.0", b: "1.1.9", want: 1},
		{a: "1.2.3", b: "1.2.4", want: -1},
		{a: "1.0.0", b: "1.0.0-rc1", want: 1}, // stable beats prerelease
		{a: "1.0.0-rc1", b: "1.0.0-rc2", want: -1},
	}
	for _, c := range cases {
		got := MustParse(c.a).Compare(MustParse(c.b))
		if got != c.want {
			t.Errorf("Compare(%s, %s) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestConstraint_Caret(t *testing.T) {
	c := MustParseConstraint("^1.2.3")

	allowed := []string{"1.2.3", "1.2.99", "1.99.0"}
	denied := []string{"1.2.2", "2.0.0", "0.9.9"}

	for _, v := range allowed {
		if !c.Allows(MustParse(v)) {
			t.Errorf("^1.2.3 should allow %s", v)
		}
	}
	for _, v := range denied {
		if c.Allows(MustParse(v)) {
			t.Errorf("^1.2.3 should DENY %s", v)
		}
	}
}

func TestConstraint_Caret_ZeroMajor(t *testing.T) {
	c := MustParseConstraint("^0.2.3")
	if !c.Allows(MustParse("0.2.99")) {
		t.Error("^0.2.3 should allow 0.2.99")
	}
	if c.Allows(MustParse("0.3.0")) {
		t.Error("^0.2.3 should DENY 0.3.0 (minor-locked when major=0)")
	}
}

func TestConstraint_Tilde(t *testing.T) {
	c := MustParseConstraint("~1.2.3")
	if !c.Allows(MustParse("1.2.99")) {
		t.Error("~1.2.3 should allow 1.2.99")
	}
	if c.Allows(MustParse("1.3.0")) {
		t.Error("~1.2.3 should DENY 1.3.0")
	}
}

func TestConstraint_Comparison(t *testing.T) {
	cases := []struct {
		c, v string
		want bool
	}{
		{c: ">=1.0.0", v: "1.0.0", want: true},
		{c: ">=1.0.0", v: "0.9.9", want: false},
		{c: ">1.0.0", v: "1.0.0", want: false},
		{c: "<2.0.0", v: "1.99.99", want: true},
		{c: "1.2.3", v: "1.2.3", want: true},
		{c: "1.2.3", v: "1.2.4", want: false},
		{c: "", v: "999.0.0", want: true},
	}
	for _, tc := range cases {
		got := MustParseConstraint(tc.c).Allows(MustParse(tc.v))
		if got != tc.want {
			t.Errorf("(%q).Allows(%q) = %v, want %v", tc.c, tc.v, got, tc.want)
		}
	}
}
