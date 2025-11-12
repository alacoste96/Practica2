package main

import "fmt"

// menú de creación
func creationMenu(option int, g *Garage) {
	switch option {
	case 1:
		g.newClient()
	case 2:
		g.newVehicle()
	case 3:
		g.newIssue()
	case 4:
		g.newMech()
	case 5:
		return
	}
}

// menu de visualización
func visualitationMenu(option int, g *Garage) {
	switch option {
	case 1:
		g.displayClients()
	case 2:
		g.displayVehicles()
	case 3:
		g.displayIssues()
	case 4:
		g.displayMechs()
	case 5:
		g.displaySlots()
	case 6:
		g.listIncFromAvehicle()
	case 7:
		g.listVFromClient()
	case 8:
		g.listDispMech()
	case 9:
		g.dispIssuesOfMech()
	case 10:
		g.listClientsVInGarage()
	case 11:
		return
	}
}

// menú de modificacion
func modificationMenu(option int, g *Garage) {
	switch option {
	case 1:
		g.modifyClient()
	case 2:
		g.modifyVehicle()
	case 3:
		g.modifyIssue()
	case 4:
		g.modifyMech()
	case 5:
		g.modifySlots()
	case 6:
		return
	}
}

// menú de eliminación
func deletionMenu(option int, g *Garage) {
	switch option {
	case 1:
		g.delClient()
	case 2:
		g.delVehicle()
	case 3:
		g.delIssue()
	case 4:
		g.delMech()
	case 5:
		return
	}
}

// menú polivalente para asignar especialidad, estados, tipos, etc
func polyAskMenuStr(prompt string, which int) string {
	var strs []string
	var option int

	switch which {
	case 1: // especialidad de mecanico
		strs = []string{string(MECHSKILL), string(ELECTRICSKILL), string(BODYSKILL)}
	case 2: // estado mecánico
		strs = []string{string(ACTIVE), string(INACTIVE)}
	case 3: // tipo incidencia
		strs = []string{string(MECHTYPE), string(ELECTRICTYPE), string(BODYTYPE)}
	case 4: // prioridad incidencia
		strs = []string{string(LOW), string(MEDIUM), string(HIGH)}
	case 5: // estado incidencia
		strs = []string{string(OPEN), string(INPROCESS), string(CLOSED)}
	}
	fmt.Printf("%s\n", prompt)
	for {
		for i, str := range strs {
			fmt.Printf("	%d. %s\n", i+1, str)
		}
		fmt.Printf("\n")
		fmt.Printf("Opción (número): ")
		fmt.Scanf("%d", &option)
		if option < 1 || option > len(strs) {
			fmt.Printf("Inválida. Pruebe de nuevo")
		} else {
			return strs[option-1]
		}
	}
}

// menú polivalente para obtener una opcion segun se requiera
func polyAskMenuInt(prompt string, which int) int {
	var strs []string
	var option int

	switch which {
	case 1: // Modificación de cliente
		strs = []string{"Cambiar ID", "Cambiar Nombre", "Cambiar teléfono", "Cambiar email"}
	case 2: // Modificación de vehículo
		strs = []string{"Cambiar matrícula", "Cambiar Marca", "Cambiar Modelo"}
	case 3: // Modificacion de incidencia
		strs = []string{"Cambiar ID", "Cambiar vehiculo asociado", "Cambiar tipo", "Cambiar prioridad",
			"Cambiar estado", "Cambiar descripción", "Cambiar mecánicos"}
	case 4: // gestionar mecánicos en una incidencia
		strs = []string{"Asignar mecánico", "Quitar mecánico"}
	case 5: // Modificación de un mecánico
		strs = []string{"Cambiar el ID", "Cambiar nombre", "Cambiar especialidad", "Cambiar el estado"}
	case 6: // asignacion de plazas de taller
		strs = []string{"Meter coche", "Sacar coche", "Volver"}
	case 7: // menú principal
		strs = []string{"Crear", "Visualizar", "Modificar", "Eliminar", "Salir"}
	case 8: // menu de creacion
		strs = []string{"Cliente", "Vehículo", "Incidencia", "Mecánico", "Volver"}
	case 9: // menu de visualizacion
		strs = []string{"Clientes", "Vehiculos", "Incidencias", "Mecánicos", "Estado actual del taller",
			"Listar Incidencias de un vehículo", "Listar todos los vehículos de un cliente",
			"Listar todos los mecánicos disponibles", "Listar incidencias asignadas a un mecánico",
			"Listar todos los clientes con vehículos en el taller", "Volver"}
	case 10: // menu de modificacion
		strs = []string{"Cliente", "Vehiculo", "Incidencia", "Mecánico",
			"Asignar vehículos a plazas del taller", "volver"}
	case 11: // menu de eliminación
		strs = []string{"Cliente", "Vehículo", "Incidencia", "Mecánico", "Volver"}
	}
	fmt.Printf("%s\n", prompt)
	for {
		for i, str := range strs {
			fmt.Printf("	%d. %s\n", i+1, str)
		}
		fmt.Printf("\n")
		fmt.Printf("Opción (número): ")
		fmt.Scanf("%d", &option)
		if option < 1 || option > len(strs) {
			fmt.Printf("Inválida. Pruebe de nuevo")
		} else {
			return option
		}
	}
}
