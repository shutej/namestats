package helpers

// http://play.golang.org/p/iLyb0E4Imh
func LevenshteinDistance(a, b *string) int {
	la := len(*a)
	la1 := la + 1
	lb := len(*b)
	lb1 := lb + 1

	d := make([][]int, la1)
	ld := len(d)
	for i := 0; i < ld; i++ {
		d[i] = make([]int, lb1)
	}
	for i := 0; i < ld; i++ {
		d[i][0] = i
	}
	ld0 := len(d[0])
	for i := 0; i < ld0; i++ {
		d[0][i] = i
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			ex := 1
			if (*a)[i-1] == (*b)[j-1] {
				ex = 0
			}
			min := d[i-1][j] + 1
			if (d[i][j-1] + 1) < min {
				min = d[i][j-1] + 1
			}
			if (d[i-1][j-1] + ex) < min {
				min = d[i-1][j-1] + ex
			}
			d[i][j] = min
		}
	}
	return d[la][lb]
}
