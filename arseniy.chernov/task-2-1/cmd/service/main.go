package main

import "fmt"

func main() {
	var n int

	_, err := fmt.Scan(&n)
	if err != nil {
		fmt.Println("invalid N")

		return
	}

	var k int
	for i := 0; i < n; i++ {
		_, err := fmt.Scan(&k)
		if err != nil {
			fmt.Println("invalid K")

			return
		}

		var t int
		var op string
		minT := 15
		maxT := 30

		for j := 0; j < k; j++ {
			_, err := fmt.Scanf("%s %d", &op, &t)
			if err != nil {
				fmt.Println("invalid operation or T")

				return
			}

			switch op {
			case ">=":
				minT = max(minT, t)
			case "<=":
				maxT = min(maxT, t)
			default:
				fmt.Println("invalid operation")

				return
			}

			if minT > maxT {
				fmt.Println(-1)
			} else {
				fmt.Println(minT)
			}
		}
	}
}
