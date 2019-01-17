# maze
maze generator in Go

```
prim := Prim{Width:5, Height:5, Seed:int64(0)}
if _, err := prim.Generate(); err != nil {
	fmt.Println(prim.String())
}
 _________
|___   ___|
|___   ___|
| |_   ___|
|_   ___  |
|_______|_|
```