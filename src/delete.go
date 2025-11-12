package main

import "fmt"

// eliminación de un cliente
func (g *Garage) delClient() {
	var id ClientID
	var c *Client

	if len(g.clients) == 0 {
		fmt.Printf("(No hay clientes que eliminar)\n\n")
		return
	}
	fmt.Printf("=== Eliminación de 1 cliente ===\n\n")
	fmt.Printf("Lista de clientes:\n")
	for _, c := range g.clients {
		fmt.Printf("	- ID: %d. Name: %s\n", c.id, c.name)
	}
	id = ClientID(askUniqueIntID("ID del cliente a eliminar: ", func(n int64) bool { return !g.ownerIDexists(ClientID(n)) }))
	fmt.Println()
	c = g.getClientByID(id)

	vehs := append([]VehicleID(nil), c.vehicles...)
	for _, vid := range vehs {
		g.delVByID(vid)
	}
	for i, c := range g.clients {
		if c.id == id {
			g.clients = append(g.clients[:i], g.clients[i+1:]...)
		}
	}
}

// Eliminación de 1 vehiculo por matricula
func (g *Garage) delVByID(id VehicleID) {
	var c *Client
	var v *Vehicle

	v = g.getVByID(id)
	for _, issueID := range v.issues {
		g.delIssueByID(issueID)
	}
	c = g.getClientByID(v.ownerID)
	for i, vid := range c.vehicles {
		if vid == id {
			c.vehicles = append(c.vehicles[:i], c.vehicles[i+1:]...)
		}
	}
	g.extractVfromSlot(id)
	for i, v := range g.vehicles {
		if v.id == id {
			g.vehicles = append(g.vehicles[:i], g.vehicles[i+1:]...)
		}
	}
}

// Eliminación de 1 vehiculo
func (g *Garage) delVehicle() {
	var id VehicleID

	if len(g.vehicles) == 0 {
		fmt.Printf("(No hay vehículos que eliminar)\n\n")
		return
	}
	fmt.Printf("=== Eliminación de un Vehículo ===\n\n")
	fmt.Printf("Lista de Matrículas:\n")
	for _, v := range g.vehicles {
		fmt.Printf("	- %s\n", v.id)
	}
	fmt.Println()
	id = VehicleID(askUniqueStrID("Matrícula: ", func(s string) bool { return !g.vidExists(VehicleID(s)) }))
	fmt.Println()
	g.delVByID(id)
}

// Eliminación de 1 incidencia
func (g *Garage) delIssue() {
	var id IncidenceID

	if len(g.issues) == 0 {
		fmt.Printf("(No hay incidencias que eliminar)\n\n")
		return
	}
	fmt.Printf("=== Eliminación de una Incidencia ===\n\n")
	fmt.Printf("Lista de IDs de incidencias:\n")
	for _, issue := range g.issues {
		fmt.Printf("	- ID: %d\n", issue.id)
	}
	fmt.Println()
	id = IncidenceID(askUniqueIntID("ID de la incidencia a eliminar: ", func(n int64) bool { return !g.issueIDexists(IncidenceID(n)) }))
	fmt.Println()
	g.delIssueByID(id)
}

// Eliminación de una Incidencia dada su ID
func (g *Garage) delIssueByID(id IncidenceID) {

	for i, issue := range g.issues {
		if issue.id == id {
			g.issues = append(g.issues[:i], g.issues[i+1:]...)
		}
	}
	for _, m := range g.mechanics {
		for i, issue := range m.issues {
			if issue == id {
				m.issues = append(m.issues[:i], m.issues[i+1:]...)
			}
		}
	}
	for _, v := range g.vehicles {
		for i, issue := range v.issues {
			if issue == id {
				v.issues = append(v.issues[:i], v.issues[i+1:]...)
			}
		}
	}
}

// Eliminación de un mecánico
func (g *Garage) delMech() {
	var id MechanicID

	if len(g.mechanics) == 0 {
		fmt.Printf("(No hay mecánicos que eliminar)\n\n")
		return
	}
	fmt.Printf("=== Eliminación de un mecánico ===\n\n")
	fmt.Printf("Lista de mecánicos:\n")
	for _, m := range g.mechanics {
		fmt.Printf("	- ID: %d. Nombre: %s\n", m.id, m.name)
	}
	id = MechanicID(askUniqueIntID("Escoja un ID: ", func(n int64) bool { return !g.mechIDexists(MechanicID(n)) }))
	fmt.Println()
	for _, issue := range g.issues {
		for i, mid := range issue.mechanics {
			if mid == id {
				issue.mechanics = append(issue.mechanics[:i], issue.mechanics[i+1:]...)
			}
		}
	}
	g.slots = g.slots[:len(g.slots)-2]
	for i, m := range g.mechanics {
		if m.id == id {
			g.mechanics = append(g.mechanics[:i], g.mechanics[i+1:]...)
		}
	}
}
