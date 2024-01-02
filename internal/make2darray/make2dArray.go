package make2darray

// Make2dArray initializes a 2d array and returns it
func Make2dArray(width, height, value int) [][]int {
	a := make([][]int, height)
	for i := range a {
		a[i] = make([]int, width)
		for j := 0; j < len(a[i]); j++ {
			a[i][j] = value
		}
	}
	return a
}
