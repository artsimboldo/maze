package maze

const (
	North = 1 << iota // 1
	South	// 2
	East	// 4
	West	// 8
)

type MazeResult [][]int

type MazeGenerator interface {
	Generate(int, int, int64) (MazeResult, error)
	String()
}