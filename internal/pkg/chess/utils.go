package chess

import (
	"fmt"
	"strings"
)

func coordToUciRowKeys() map[int]byte {
	return map[int]byte{0: '1', 1: '2', 2: '3', 3: '4', 4: '5', 5: '6', 6: '7', 7: '8'}
}

func coordToUciColumnKeys() map[int]byte {
	return map[int]byte{0: 'a', 1: 'b', 2: 'c', 3: 'd', 4: 'e', 5: 'f', 6: 'g', 7: 'h'}
}

func mustCoordsToUcis(from, to *coord) string {
	if str, err := coordsToUcis(from, to); err != nil {
		panic(err)
	} else {
		return str
	}
}

func coordsToUcis(from, to *coord) (string, error) {
	resultBuilder := &strings.Builder{}

	coordToUciRowKeys := coordToUciRowKeys()
	coordToUciColumnKeys := coordToUciColumnKeys()

	fromCByte, exists := coordToUciColumnKeys[from.c]
	if exists == false {
		return "", fmt.Errorf("coord %c does not exist", from.c)
	}

	fromRByte, exists := coordToUciRowKeys[from.r]
	if exists == false {
		return "", fmt.Errorf("coord %c does not exist", from.r)
	}

	toCByte, exists := coordToUciColumnKeys[to.c]
	if exists == false {
		return "", fmt.Errorf("coord %c does not exist", to.c)
	}

	toRByte, exists := coordToUciRowKeys[to.r]
	if exists == false {
		return "", fmt.Errorf("coord %c does not exist", to.r)
	}

	fmt.Fprintf(resultBuilder, "%c%c%c%c", fromCByte, fromRByte, toCByte, toRByte)

	return resultBuilder.String(), nil
}
