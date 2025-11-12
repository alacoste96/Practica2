package main

import (
	"fmt"
	"strings"
)

// listar todas las incidencias en el taller
func (g *Garage) listIssuesInGarage() {
	if len(g.issues) == 0 {
		fmt.Printf("(No hay incidencias en el taller)\n\n")
		return
	}
	fmt.Printf("=== Lista de todas las incidencias en el taller ===\n\n")
	for _, s := range g.slots {
		v := g.getVByID(*s.vehicleID)
		for _, isueID := range v.issues {
			i := g.getIssueByID(isueID)
			g.dispAnIssue(i)
		}
	}
}

// listar los clientes con vehículos en el taller
func (g *Garage) listClientsVInGarage() {

	var noclients = true

	if len(g.clients) == 0 {
		fmt.Printf("(No hay clientes registrados)\n\n")
		return
	}
	fmt.Printf("=== Lista de Clientes con vehículos en el taller ===\n\n")
	for _, c := range g.clients {
		if g.hasVInSlot(c) {
			g.dispAclient(c)
			noclients = false
		}
	}
	if noclients {
		fmt.Printf("	(No hay clientes con vehículos en el taller)\n")
	}
}

// Visualizar taller (plazas ocupadas/libres)
func (g *Garage) displaySlots() {
	if len(g.slots) == 0 {
		fmt.Printf("(No hay mecánicos ni plazas en el taller)\n\n")
		return
	}
	fmt.Printf("=== Estado del Taller ===\n\n")
	for _, s := range g.slots {
		fmt.Printf("Plaza %d: ", s.number)
		if s.vehicleID == nil {
			fmt.Printf("Libre\n")
		} else {
			fmt.Printf("Ocupada\n")
		}
	}
	fmt.Println()
}

// Visualizar 1 cliente
func (g *Garage) dispAclient(c *Client) {
	fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos
	fmt.Printf("# Cliente %d\n", c.id)
	fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos

	// Datos principales
	fmt.Printf("Nombre:   %s\n", c.name)
	fmt.Printf("Teléfono: %s\n", c.phone)
	fmt.Printf("Email:    %s\n", c.email)
	fmt.Printf("Coches asociados:\n")
	for _, v := range c.vehicles {
		fmt.Printf("	- %s\n", v)
	}
	fmt.Println()
}

// Visualizar clientes
func (g *Garage) displayClients() {

	if len(g.clients) == 0 {
		fmt.Printf("(No hay clientes registrados)\n\n")
		return
	}
	fmt.Printf("=== Lista de Clientes ===\n\n")
	for _, c := range g.clients {
		g.dispAclient(c)
	}
	fmt.Println()
}

// Visualizar vehículos
func (g *Garage) displayVehicles() {

	if len(g.vehicles) == 0 {
		fmt.Printf("(No hay vehículos registrados)\n\n")
		return
	}
	fmt.Printf("=== Lista de Vehículos ===\n\n")
	for _, v := range g.vehicles {
		// Cabecera
		owner := g.getClientByID(v.ownerID)
		fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos
		fmt.Printf("# Vehículo %s\n", v.id)
		fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos

		// Datos principales
		fmt.Printf("Marca:           %s\n", v.brand)
		fmt.Printf("Modelo:          %s\n", v.model)
		fmt.Printf("Dueño:           %s\n", owner.name)
		if v.checkInAt.IsZero() {
			fmt.Printf("Entrada:         (Por asignar)\n")
		} else {
			fmt.Printf("Entrada:         %s\n", v.checkInAt.Format("02/01/2006 15:04"))
		}
		if v.eta.IsZero() {
			fmt.Printf("Salida estimada: (Por asignar)\n")
		} else {
			fmt.Printf("Salida estimada: %s\n", v.eta.Format("02/01/2006 15:04"))
		}
		// Incidencias en bloque
		fmt.Printf("Incidencias:\n")
		if len(v.issues) == 0 {
			fmt.Printf("	(Sin Incidencias)")
		} else {
			for _, issueID := range v.issues {
				//fmt.Printf("ID: %d\n", issueID)
				issue := g.getIssueByID(issueID)
				fmt.Printf("	- ID: %d. ", issue.id)
				fmt.Printf("Descripción: %s\n", issue.description)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Visualizar 1 Incidencia
func (g *Garage) dispAnIssue(i *Incidence) {
	// Cabecera
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("# Incidencia %d\n", i.id)
	fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos

	// Datos principales
	fmt.Printf("Vehículo:   %s\n", i.vehicleID)
	fmt.Printf("Tipo:       %s\n", i.kind)
	fmt.Printf("Prioridad:  %s\n", i.prio)
	fmt.Printf("Estado:     %s\n", i.status)

	// Descripción (en bloque)
	fmt.Println("Descripción:")
	fmt.Printf("  %s\n", i.description)

	// Mecánicos asignados
	if len(i.mechanics) == 0 {
		fmt.Printf("Mecánicos:\n	(Sin mecánicos asignados)\n")
	} else {
		fmt.Println("Mecánicos asignados:")
		for _, mechID := range i.mechanics {
			m := g.getMechByID(mechID)
			name := "(desconocido)"
			if m != nil && m.name != "" {
				name = m.name
			}
			fmt.Printf("  - %s\n", name)
		}
	}
	fmt.Println()
}

// Visualizar Incidencias
func (g *Garage) displayIssues() {

	if len(g.issues) == 0 {
		fmt.Printf("(no hay incidencias)\n\n")
		return
	}
	fmt.Println("=== Lista de Incidencias ===")
	for _, i := range g.issues {
		g.dispAnIssue(i)
	}
	fmt.Println()
}

// listar Incidencias asignadas a un mecánico
func (g *Garage) dispIssuesOfMech() {
	var mech *Mechanic
	var mid MechanicID

	if len(g.mechanics) == 0 {
		fmt.Printf("(No hay mecánicos registrados)\n\n")
		return
	}
	fmt.Printf("Lista de mecánicos:\n")
	for _, m := range g.mechanics {
		fmt.Printf("	- ID: %d. Name: %s\n", m.id, m.name)
	}
	mid = MechanicID(askUniqueIntID("Introduzca el ID del mecánico: ", func(n int64) bool { return !g.mechIDexists(MechanicID(n)) }))
	fmt.Println()
	mech = g.getMechByID(mid)
	fmt.Printf("=== Incidencias asignadas a %s ===\n\n", mech.name)
	for _, isueID := range mech.issues {
		i := g.getIssueByID(isueID)
		g.dispAnIssue(i)
	}
}

// Visualizar 1 mecánico
func (g *Garage) displayMech(m *Mechanic) {
	// Cabecera
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("# Mecánico %d\n", m.id)
	fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos

	// Datos principales
	fmt.Printf("Nombre:       %s\n", m.name)
	fmt.Printf("Especialidad: %s\n", m.skill)
	fmt.Printf("Experiencia:  %d años\n", m.experience)
	fmt.Printf("Estado:       %s\n", m.status)
	fmt.Println()
}

// Visualizar Mecánicos
func (g *Garage) displayMechs() {

	if len(g.mechanics) == 0 {
		fmt.Printf("(No hay mecánicos registrados)\n\n")
		return
	}
	fmt.Printf("=== Lista de Mecánicos ===\n\n")
	for _, m := range g.mechanics {
		g.displayMech(m)
	}
	fmt.Println()
}

// Visualizar incidencias de un vehículo
func (g *Garage) listIncFromAvehicle() {
	var vid VehicleID

	if len(g.vehicles) == 0 {
		fmt.Printf("(No hay vehículos registrados)\n\n")
		return
	}
	fmt.Printf("=== Listar Incidencias de 1 Vehículo ===\n\n")
	fmt.Printf("Lista de vehículos:\n")
	for _, v := range g.vehicles {
		fmt.Printf("	- %s\n", v.id)
	}
	fmt.Println()
	vid = VehicleID(askUniqueStrID("Escriba la matrícula del vehículo: ", func(s string) bool { return !g.vidExists(VehicleID(s)) }))
	fmt.Println()
	fmt.Printf("=== Incidencias del vehículo %s ===\n\n", vid)
	for _, issueID := range g.getVByID(vid).issues {
		i := g.getIssueByID(issueID)
		g.dispAnIssue(i)
	}
}

// Visualizar vehículos de un cliente
func (g *Garage) listVFromClient() {
	var client *Client
	var cid ClientID

	if len(g.clients) == 0 {
		fmt.Printf("(No hay clientes registrados)\n\n")
		return
	}
	fmt.Printf("=== Listar Vehículos de un Cliente ===\n\n")
	fmt.Printf("Lista de clientes:\n")
	for _, c := range g.clients {
		fmt.Printf("	- ID: %d. Nombre: %s\n", c.id, c.name)
	}
	fmt.Println()
	cid = ClientID(askUniqueIntID("Introduzca el ID del cliente: ", func(n int64) bool { return !g.ownerIDexists(ClientID(n)) }))
	fmt.Println()
	client = g.getClientByID(cid)
	fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos
	fmt.Printf("Vehículos de %s:\n", client.name)
	fmt.Println(strings.Repeat("-", 60)) // 60 guiones seguidos
	if len(client.vehicles) == 0 {
		fmt.Printf("(Sin vehículos asociados)\n\n")
		return
	}
	for _, vid := range client.vehicles {
		v := g.getVByID(vid)
		fmt.Printf("	- Matrícula: %s\n", vid)
		fmt.Printf("		Marca: %s\n", v.brand)
		fmt.Printf("		Modelo: %s\n", v.model)
	}
	fmt.Println()
}

// listar todos los mecáncos disponibles (no asignados a incidencias)
func (g *Garage) listDispMech() {
	var ndisp int = 0

	if len(g.mechanics) == 0 {
		fmt.Printf("(No hay mecánicos registrados)\n\n")
		return
	}
	fmt.Printf("=== Mecánicos Disponibles ===\n\n")
	for _, m := range g.mechanics {
		if len(m.issues) == 0 {
			g.displayMech(m)
			ndisp++
		}
	}
	if ndisp == 0 {
		fmt.Printf("\n(No hay mecánicos disponibles)\n")
	}
}
