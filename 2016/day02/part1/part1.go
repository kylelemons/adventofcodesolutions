package main

import (
	"bufio"
	"fmt"
	"strings"
)

var pad = [][]string{
	{"1", "2", "3"},
	{"4", "5", "6"},
	{"7", "8", "9"},
}

func code(input string) string {
	x, y := 1, 1
	move := func(newx, newy int) {
		switch {
		case newx < 0:
		case newx >= len(pad[0]):
		case newy < 0:
		case newy >= len(pad):
		default:
			x, y = newx, newy
		}
	}

	var code []string
	lines := bufio.NewScanner(strings.NewReader(input))
	for lines.Scan() {
		for _, c := range lines.Text() {
			switch c {
			case 'U':
				move(x, y-1)
			case 'D':
				move(x, y+1)
			case 'L':
				move(x-1, y)
			case 'R':
				move(x+1, y)
			}
		}
		code = append(code, pad[y][x])
	}
	return strings.Join(code, "")
}

func main() {
	for i, input := range []string{
		`ULL
RRDDD
LURDL
UUUUD`,
		`URULLLLLRLDDUURRRULLLDURRDRDRDLURURURLDLLLLRUDDRRLUDDDDDDLRLRDDDUUDUDLDULUDLDURDULLRDDURLLLRRRLLRURLLUDRDLLRRLDDRUDULRRDDLUUUDRLDLURRRULURRDLLLDDDLUDURDDRLDDDLLRULDRUDDDLUDLURUDLLRURRUURUDLLLUUUUDDURDRDDDLDRRUDURDLLLULUDURURDUUULRULUDRUUUUDLRLUUUUUDDRRDDDURULLLRRLDURLDLDRDLLLUULLRRLLLLDRLRDRRDRRUDDLULUUDDDDRRUUDDLURLRDUUDRRLDUDLRRRLRRUUDURDRULULRDURDRRRDLDUUULRDDLRLRDLUUDDUDDRLRRULLLULULLDDDRRDUUUDDRURDDURDRLRDLDRDRULRLUURUDRLULRLURLRRULDRLRDUDLDURLLRLUDLUDDURDUURLUDRLUL
LLLUUURUULDDDULRRDLRLLLLLLLLRURRDLURLUDRRDDULDRRRRRRLDURRULDDULLDDDRUUDLUDULLDLRRLUULULRULURDURLLDULURDUDLRRLRLLDULLRLDURRUULDLDULLRDULULLLULDRLDLDLDLDDLULRLDUDRULUDDRDDRLRLURURRDULLUULLDRRDRRDLDLLRDLDDUUURLUULDDRRRUULDULDDRDDLULUDRURUULLUDRURDRULDRUULLRRDURUDDLDUULLDDRLRRDUDRLRRRLDRLRULDRDRRUDRLLLDDUDLULLURRURRLUURDRLLDLLDUDLUUURRLRDDUDRLUDLLRULLDUUURDLUUUDUDULRLDLDRUUDULRDRRUDLULRLRDLDRRDDDUDLDLDLRUURLDLLUURDLDLRDLDRUDDUURLLLRDRDRRULLRLRDULUDDDLUDURLDUDLLRULRDURDRDLLULRRDLLLDUURRDUDDLDDRULRRRRLRDDRURLLRRLLL
DRURLDDDDRLUDRDURUDDULLRRLLRLDDRLULURLDURRLDRRLRLUURDDRRDLRDLDLULDURUDRLRUDULRURURLRUDRLLDDUDDRDLDRLLDDLRRDRUUULDUUDRUULRLLDLLULLLRRDRURDLDDRRDDUDDULLDUUULDRUDLDLURLDRURUDLRDDDURRLRDDUDLLLRRUDRULRULRRLLUUULDRLRRRLLLDLLDUDDUUDRURLDLRRUUURLUDDDRRDDLDDDDLUURDDULDRLRURLULLURRDRLLURLLLURDURLDLUDUUDUULLRLDLLLLULRDDLDUDUDDDUULURRLULDLDRLRDRLULLUDDUUUUURDRURLDUULDRRDULUDUDLDDRDLUDDURUDURLDULRUDRRDLRLRDRRURLDLURLULULDDUUDLRLLLLURRURULDDRUUULLDULDRDULDDDLLLRLULDDUDLRUDUDUDURLURLDDLRULDLURD
DRUDRDURUURDLRLUUUUURUDLRDUURLLDUULDUULDLURDDUULDRDDRDULUDDDRRRRLDDUURLRDLLRLRURDRRRDURDULRLDRDURUDLLDDULRDUDULRRLLUDLLUUURDULRDDLURULRURDDLRLLULUDURDRRUDLULLRLDUDLURUDRUULDUDLRDUDRRDULDDLDRLRRULURULUURDULRRLDLDULULRUUUUULUURLURLRDLLRRRRLURRUDLRLDDDLDRDRURLULRDUDLRLURRDRRLRLLDLDDLLRRULRLRLRUDRUUULLDUULLDDRLUDDRURLRLDLULDURLLRRLDLLRDDDUDDUULLUDRUDURLLRDRUDLUDLLUDRUUDLRUURRRLLUULLUUURLLLRURUULLDLLDURUUUULDDDLRLURDRLRRRRRRUDLLLRUUULDRRDLRDLLDRDLDDLDLRDUDLDDRDDDDRULRRLRDULLDULULULRULLRRLLUURUUUDLDLUDUDDDLUUDDDDUDDDUURUUDRDURRLUULRRDUUDDUDRRRDLRDRLDLRRURUUDRRRUUDLDRLRDURD
DDDLRURUDRRRURUUDLRLRDULDRDUULRURRRUULUDULDDLRRLLRLDDLURLRUDRLRRLRDLRLLDDLULDLRRURDDRDLLDDRUDRRRURRDUDULUDDULRRDRLDUULDLLLDRLUDRDURDRRDLLLLRRLRLLULRURUUDDRULDLLRULDRDLUDLULDDDLLUULRRLDDUURDLULUULULRDDDLDUDDLLLRRLLLDULRDDLRRUDDRDDLLLLDLDLULRRRDUDURRLUUDLLLLDUUULDULRDRULLRDRUDULRUUDULULDRDLDUDRRLRRDRLDUDLULLUDDLURLUUUDRDUDRULULDRDLRDRRLDDRRLUURDRULDLRRLLRRLDLRRLDLDRULDDRLURDULRRUDURRUURDUUURULUUUDLRRLDRDLULDURUDUDLUDDDULULRULDRRRLRURLRLRLUDDLUUDRRRLUUUDURLDRLRRDRRDURLLL`,
	} {
		fmt.Printf("%d. code(%q)\n = %s\n", i, input, code(input))
	}
}
