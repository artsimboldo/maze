package maze

import (
	"math/rand"
	"bytes"
	"strings"
	"errors"
)

/*
PRISM maze generaator
see http://weblog.jamisbuck.org/2011/1/10/maze-generation-prim-s-algorithm.html
*/

const (
	_in_ = 0x10
	_frontier_ = 0x20
)

var opposite = map[int]int{East:West, West:East, North:South, South:North}

type Prim struct {
	Width, Height int
	Seed int64
	grid MazeResult
}

type pair struct {
	x, y int
}

/*
addFrontier() adds a cell as part of the frontier to explore
*/
func (pr *Prim) addFrontier(x, y int, frontier *[]pair) {
	if x >= 0 && y >= 0 && y < pr.Height && x < pr.Width && pr.grid[y][x] == 0 {
		pr.grid[y][x] |= _frontier_
		*frontier = append(*frontier, pair{x, y})
	}
}

/*
mark() marks a cell as part (in) of the maze solution 
*/
func (pr *Prim) mark(x, y int, frontier *[]pair) {
	pr.grid[y][x] &= ^_frontier_
	pr.grid[y][x] |= _in_
	pr.addFrontier(x - 1, y, frontier)
	pr.addFrontier(x + 1, y, frontier)
	pr.addFrontier(x, y - 1, frontier)
	pr.addFrontier(x, y + 1, frontier)
}

/*
neighbors() appends possible cell's neighborhood not already in the maze solution 
*/
func (pr Prim) neighbors(x, y int) ([]pair, int) {
	var n []pair

	// neighbors are cell not "in"
	if x > 0 && pr.grid[y][x-1] & _in_ != 0 {
		n = append(n, pair{x-1,y})
	}

	if x + 1 < pr.Width && pr.grid[y][x+1] & _in_ != 0 {
		n = append(n, pair{x+1, y})
	}

	if y > 0 && pr.grid[y-1][x] & _in_ != 0 {
		n = append(n, pair{x, y-1})
	}

	if y + 1 < pr.Height && pr.grid[y+1][x] & _in_ != 0 {
		n = append(n, pair{x, y+1})
	}

	return n, len(n)
}

/*
direction() returns cardinal directions reative to two pairs of coordinates
*/
func direction(f, t pair) int {
	switch {
		case f.x < t.x: return East
		case f.x > t.x: return West
		case f.y < t.y: return South
//		case f.y > t.y:
		default: return North
	}
}

/*
String() displays the maze as a string
*/
func (pr Prim) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(" " + strings.Repeat("_", pr.Width * 2 - 1) + "\n")
	for _, row := range pr.grid {
		buffer.WriteString("|")
		for x, cell := range row {
		    if cell & South != 0 {
	     		buffer.WriteString(" ")
			} else {
				// add south wall
	     		buffer.WriteString("_")
			}
			if cell & East != 0 {
      			// specific rule to fill south holes due to east walls shifting
				if (cell | row[x+1]) & South != 0 {
					buffer.WriteString(" ")
				} else {
					buffer.WriteString("_")
				}
			} else {
				// add east wall
	     		buffer.WriteString("|")
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

/*
Generate() is the main routine to build the maze, returns the solution in MazeResult
*/
func (pr *Prim) Generate() (MazeResult, error) {
	if pr.Width < 2 || pr.Height < 2 {
		return nil, errors.New("Error: size of width and height must be >1 !")
	}

	var (
		p, q pair
		frontier []pair
		grid = make(MazeResult, pr.Height)
	)

	// create an width x height empty grid
	for i := range grid {
		grid[i] = make([]int, pr.Width)
	}
	pr.grid = grid

	// choose at random first cell and compute initial frontier
	rand.Seed(pr.Seed)
	pr.mark(
		rand.Intn(pr.Width - 1), 
		rand.Intn(pr.Height - 1), 
		&frontier)

	for {
		if sz := len(frontier); sz == 0 {
			// we're done
			break
		// retrieve cell p from the frontier
		} else if sz == 1 {
			p = frontier[0]
			frontier = frontier[:0]
		} else {
			i := rand.Intn(sz - 1)
			p = frontier[i]
			frontier = append(frontier[:i], frontier[i+1:]...)
		}
		// retrieve neighbors of p and select q
		if n, sz := pr.neighbors(p.x, p.y); sz > 1 {
			q = n[rand.Intn(sz - 1)]
		} else {
			q = n[0]
		}
		// record path between p and q
		dir := direction(p, q)
		grid[p.y][p.x] |= dir
		grid[q.y][q.x] |= opposite[dir]
		// mark p as "_in_" and compute new frontier
		pr.mark(p.x, p.y, &frontier)
	}
	return grid, nil
}