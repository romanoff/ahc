package preprocessor

type Renamer struct {
	Counter int
	Classes map[string]string
}

var validClassSymbols = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func GetMinifiedCssClass(number int) string {
	divider := 1
	numberBytes := []byte{}
	base := len(validClassSymbols)
	for {
		divider = divider * base
		byteNum := ((number % divider) - (number % (divider / base))) / (divider / base)
		if byteNum == 0 {
			if number == 0 {
				numberBytes = append(numberBytes, validClassSymbols[0])
			}
			break
		}
		numberBytes = append(numberBytes, validClassSymbols[byteNum])
	}
	return string(numberBytes)
}
