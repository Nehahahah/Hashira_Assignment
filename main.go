package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
)

type Root struct {
	Base  string `json:"base"`
	Value string `json:"value"`
}

type Input struct {
	Keys int    `json:"keys"`
	Roots []Root `json:"roots"`
}

func decodeValue(baseStr, value string) *big.Int {
	base := 10
	if baseStr == "10" {
		base = 10
	} else if baseStr == "2" {
		base = 2
	} else if baseStr == "8" {
		base = 8
	} else if baseStr == "16" {
		base = 16
	}
	n := new(big.Int)
	n.SetString(value, base)
	return n
}

func lagrangeInterpolation(points [][2]*big.Int) *big.Int {
	result := big.NewInt(0)
	for i := 0; i < len(points); i++ {
		term := new(big.Int).Set(points[i][1]) // yi
		for j := 0; j < len(points); j++ {
			if i == j {
				continue
			}
			num := new(big.Int).Neg(points[j][0])                // -xj
			den := new(big.Int).Sub(points[i][0], points[j][0]) // xi - xj
			frac := new(big.Rat).SetFrac(num, den)
			termRat := new(big.Rat).Mul(new(big.Rat).SetInt(term), frac)
			term.Quo(termRat.Num(), termRat.Denom())
		}
		result.Add(result, term)
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <inputfile.json>")
		return
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	var input Input
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	points := make([][2]*big.Int, len(input.Roots))
	for i, root := range input.Roots {
		points[i][0] = big.NewInt(int64(i + 1))
		points[i][1] = decodeValue(root.Base, root.Value)
	}

	C := lagrangeInterpolation(points[:input.Keys])
	fmt.Println("Secret:", C)
}
