package main

import (
	"bufio"
	"fmt"
	"os"
)

// contar plazas libres
func (g *Garage) availableSlots() int {
	var count int = 0

	for _, s := range g.slots {
		if s.vehicleID == nil {
			count++
		}
	}
	return count
}

// devuelve si un vehículo está en una plaza del taller
func (g *Garage) vIsInSlot(vid VehicleID) bool {
	for _, v := range g.slots {
		if v.vehicleID != nil && *v.vehicleID == vid {
			return true
		}
	}
	return false
}

// devuelve si un mecanicos está sociado a una incidencia o no
func isInIncidence(i *Incidence, m MechanicID) bool {
	for _, mid := range i.mechanics {
		if mid == m {
			return true
		}
	}
	return false
}

// eliminar mecánicos de una incidencia
func (g *Garage) removeMechsFromIssue(issue *Incidence) {
	var midtodel MechanicID
	var m *Mechanic

	if len(issue.mechanics) == 0 {
		fmt.Printf("(La incidencia no tiene mecanicos asignados)\n\n")
		return
	}
	fmt.Printf("Lista de mecánicos de la incidencia.\n")
	for _, mid := range issue.mechanics {
		m := g.getMechByID(mid)
		fmt.Printf("	- ID: %d. Nombre: %s\n", m.id, m.name)
	}
	fmt.Println()
	midtodel = MechanicID(askUniqueIntID("Introduzca el ID del mecánico a quitar: ", func(n int64) bool { return !g.mechIDexists(MechanicID(n)) }))
	for idx, mid := range issue.mechanics {
		if mid == midtodel {
			issue.mechanics = append(issue.mechanics[:idx], issue.mechanics[idx+1:]...)
		}
	}
	m = g.getMechByID(midtodel)
	for idx, issueID := range m.issues {
		if issueID == issue.id {
			m.issues = append(m.issues[:idx], m.issues[idx+1:]...)
		}
	}
}

// asignar mecánicos a una incidencia
func (g *Garage) assignMechsToIssue(i *Incidence) {
	var mcomp []Mechanic
	var option int
	var mech *Mechanic

	for _, m := range g.mechanics {
		if string(i.kind) == string(m.skill) && !isInIncidence(i, m.id) && m.status == ACTIVE {
			mcomp = append(mcomp, *m)
		}
	}
	if len(mcomp) == 0 {
		fmt.Printf("No hay mecánicos compatibles. Contrate uno especializado en %s.\n", i.kind)
		return
	}
	fmt.Printf("Asignación de Mecánicos:\n\n")
	for {
		fmt.Printf("Mecánicos especializados en %s:\n", i.kind)
		for i, m := range mcomp {
			fmt.Printf("	%d. Nombre: %s. ID: %d\n", i+1, m.name, m.id)
		}
		fmt.Printf("	%d. No asignar más\n", len(mcomp)+1)
		fmt.Printf("\nEscoja uno(rango 1 - %d): ", len(mcomp)+1)
		fmt.Scanf("%d", &option)
		if option < 1 || option > (len(mcomp)+1) {
			fmt.Printf("Opción inválida. Pruebe de nuevo.\n")
		} else if option == (len(mcomp) + 1) {
			return
		} else {
			mech = g.getMechByID(mcomp[option-1].id)
			mech.issues = append(mech.issues, i.id)
			i.mechanics = append(i.mechanics, mcomp[option-1].id)
			mcomp = append(mcomp[:option-1], mcomp[option:]...)
			if len(mcomp) == 0 {
				fmt.Printf("No quedan más mecánicos compatibles.\n")
				return
			}
		}
		fmt.Println()
	}
}

// pedir una string normal
func askStr(prompt string) string {
	fmt.Printf("%s", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// pedir un ID único de tipo entero54
func askUniqueIntID(prompt string, exists func(int64) bool) int64 {
	var id int64
	for {
		fmt.Print(prompt)
		fmt.Scanf("%d", &id) //fmt.Scanln()
		if id <= 0 {
			fmt.Println("Debe ser > 0.")
			continue
		}
		if exists(id) {
			fmt.Println("Incorrecto, pruebe otro.")
			continue
		}
		return id
	}
}

// pedir una matrícula única de tipo string
func askUniqueStrID(prompt string, exists func(string) bool) string {
	var s string
	for {
		fmt.Print(prompt)
		fmt.Scanf("%s", &s) //fmt.Scanln()
		if s == "" {
			fmt.Println("No puede estar vacío.")
			continue
		}
		if exists(s) {
			fmt.Println("Incorrecta, pruebe otro.")
			continue
		}
		return s
	}
}

// OBTENER si una incidencia existe
func (g *Garage) issueIDexists(id IncidenceID) bool {
	for _, i := range g.issues {
		if i.id == id {
			return true
		}
	}
	return false
}

// OBTENER si una MATRICULA existe
func (g *Garage) vidExists(vid VehicleID) bool {
	for _, v := range g.vehicles {
		if v.id == vid {
			return true
		}
	}
	return false
}

// OBTENER si un ID de un cliente existe
func (g *Garage) ownerIDexists(id ClientID) bool {
	for _, c := range g.clients {
		if c.id == id {
			return true
		}
	}
	return false
}

// OBTENER si un ID de un mecánico existe
func (g *Garage) mechIDexists(id MechanicID) bool {
	for _, m := range g.mechanics {
		if m.id == id {
			return true
		}
	}
	return false
}

// OBTENER si un VEHICULO de un CLIENTE está en el taller
func (g *Garage) hasVInSlot(c *Client) bool {
	for _, vid := range c.vehicles {
		for _, s := range g.slots {
			if s.vehicleID != nil && *s.vehicleID == vid {
				return true
			}
		}
	}
	return false
}

// OBTENER un * a CLIENTE a partir de su CLAVE
func (g *Garage) getClientByID(id ClientID) *Client {
	for _, c := range g.clients {
		if c.id == id {
			return c
		}
	}
	return nil
}

// OBTENER un * a MECÁNICO a partir de su CLAVE
func (g *Garage) getMechByID(id MechanicID) *Mechanic {
	for _, m := range g.mechanics {
		if m.id == id {
			return m
		}
	}
	return nil
}

// OBTENER un * a INCIDENCIA a partir de su CLAVE
func (g *Garage) getIssueByID(id IncidenceID) *Incidence {
	for _, i := range g.issues {
		if i.id == id {
			return i
		}
	}
	return nil
}

// OBTENER un * a VEHICULO a partir de su CLAVE
func (g *Garage) getVByID(vid VehicleID) *Vehicle {
	for _, v := range g.vehicles {
		if v.id == vid {
			return v
		}
	}
	return nil
}
