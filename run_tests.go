package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func main() {

	files, _ := filepath.Glob("teste_*.txt")
	sort.Strings(files)

	for _, file := range files {
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Testando arquivo: %s\n", file)

		content, _ := os.ReadFile(file)
		fmt.Printf("Entrada: %s\n", string(content))
		fmt.Println(">>> Saída:")

		cmd := exec.Command("go", "run", "analise_lexica.go", file)
		output, _ := cmd.CombinedOutput()
		fmt.Print(string(output))
	}
}
