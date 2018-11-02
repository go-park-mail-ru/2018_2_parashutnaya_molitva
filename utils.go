package chess

import "fmt"

func UciToCoordRowKeys() map[byte]int {
	rows := map[byte]int{
		'1': 0, '2': 1, '3': 2, '4': 3, '5': 4, '6': 5, '7': 6, '8': 7,
	}
	return rows
}

func UciToCoordColumnKeys() map[byte]int {
	columns := map[byte]int{
		'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5, 'g': 6, 'h': 7,
	}
	return columns
}

func CoordToUciRowKeys() map[int]byte {
	rows := map[int]byte{
		0: '1', 1: '2', 2: '3', 3: '4', 4: '5', 5: '6', 6: '7', 7: '8',
	}
	return rows
}

func CoordToUciColumnKeys() map[int]byte {
	columns := map[int]byte{
		0: 'a', 1: 'b', 2: 'c', 3: 'd', 4: 'e', 5: 'f', 6: 'g', 7: 'h',
	}
	return columns
}

func UcisToCoords(uci string) (from, to Coord) {
	if len(uci) != 4 {
		panic(fmt.Sprintf("uci string must be 4 characters long (%s)", uci))
	}

	uciToCoordColumnKeys := UciToCoordColumnKeys()
	uciToCoordRowKeys := UciToCoordRowKeys()

	val, exists := uciToCoordColumnKeys[uci[0]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[0]))
	}
	from.c = val

	val, exists = uciToCoordRowKeys[uci[1]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[1]))
	}
	from.r = val

	val, exists = uciToCoordColumnKeys[uci[2]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[2]))
	}
	to.c = val

	val, exists = uciToCoordRowKeys[uci[3]]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", uci[3]))
	}
	to.r = val

	return
}

func CoordsToUcis(from, to Coord) string {
	var result string

	coordToUciRowKeys := CoordToUciRowKeys()
	coordToUciColumnKeys := CoordToUciColumnKeys()

	val, exists := coordToUciColumnKeys[from.c]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", from.c))
	}
	result += string(val)

	val, exists = coordToUciRowKeys[from.r]
	if exists == false {
		panic(fmt.Sprintf("coord %d does not exist", from.r))
	}
	result += string(val)

	val, exists = coordToUciColumnKeys[to.c]
	if exists == false {
		panic(fmt.Sprintf("coord %c does not exist", to.c))
	}
	result += string(val)

	val, exists = coordToUciRowKeys[to.r]
	if exists == false {
		panic(fmt.Sprintf("coord %d does not exist", to.r))
	}
	result += string(val)

	return result
}
