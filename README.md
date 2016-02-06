# go package extreme


A simple extreme values generator that can be used in combination
 with testing/quick.
 
Example:
 
 func TestAbs2e(t *testing.T) {
	f := func(x int) bool {
		return Abs2(x) >= 0       // post-condition of Abs2
	}
	if err := quick.Check(f, &quick.Config{Values: extremeValues.ExtremeValues(f)}); err != nil {
		t.Error(err)
	}
