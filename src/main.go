package main

import "fmt"

func main() {
	var option int
	var g *Garage

	g = newGarage()
	fmt.Printf("-----------------------------------\n")
	fmt.Printf("- BIENVENIDO AL TALLER LACOSTE.CO -\n")
	fmt.Printf("-----------------------------------\n")
	for {
		fmt.Printf("=== Menú Principal ===\n\n")
		option = polyAskMenuInt("Seleccione una opción:", 7)
		switch option {
		case 1:
			fmt.Printf("=== Menú de Creación ===\n\n")
			option = polyAskMenuInt("Seleccione qué quiere crear: ", 8)
			creationMenu(option, g)
		case 2:
			fmt.Printf("=== Menú de Visualización ===\n\n")
			option = polyAskMenuInt("Seleccione qué visualizar: ", 9)
			fmt.Println()
			visualitationMenu(option, g)
		case 3:
			fmt.Printf("=== Menú de Modificación ===\n\n")
			option = polyAskMenuInt("Seleccione qué modificar: ", 10)
			fmt.Println()
			modificationMenu(option, g)
		case 4:
			fmt.Printf("=== Menú de Eliminación ===\n\n")
			option = polyAskMenuInt("Seleccione qué eliminar: ", 11)
			fmt.Println()
			deletionMenu(option, g)
		case 5:
			return
		}
	}
}
