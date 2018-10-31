package chess

import "fmt"

func UciToCoords(uci string) (from, to Coord) {
	if len(uci) != 4 {
		panic(fmt.Sprintf("uci string must be 4 characters long (%s)", uci))
	}
	rows := map[byte]int{
		'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5, 'g': 6, 'h': 7,
		'1': 0, '2': 1, '3': 2, '4': 3, '5': 4, '6': 5, '7': 6, '8': 7,
	}

	val, exists := rows[uci[0]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[0]))
	}
	from.c = val

	val, exists = rows[uci[1]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[1]))
	}
	from.r = val

	val, exists = rows[uci[2]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[2]))
	}
	to.c = val

	val, exists = rows[uci[3]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[3]))
	}
	to.r = val

	return
}
