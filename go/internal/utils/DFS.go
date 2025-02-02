package utils

// FindContoursDFS apply DFS to find contours of A4
func FindContoursDFS(image [][]int) [][][]int {
	height, width := len(image), len(image[0])
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	var contours [][][]int
	directions := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if image[y][x] == 1 && !visited[y][x] {
				contour := [][]int{}
				stack := [][]int{{y, x}}

				for len(stack) > 0 {
					point := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					py, px := point[0], point[1]

					if visited[py][px] {
						continue
					}
					visited[py][px] = true
					contour = append(contour, []int{px, py})

					for _, d := range directions {
						ny, nx := py+d[0], px+d[1]
						if ny >= 0 && ny < height && nx >= 0 && nx < width && image[ny][nx] == 1 && !visited[ny][nx] {
							stack = append(stack, []int{ny, nx})
						}
					}
				}
				contours = append(contours, contour)
			}
		}
	}
	return contours
}
