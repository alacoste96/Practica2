package main

import "fmt"

// modificar el id de un cliente
func (g *Garage) modifyClientID(c *Client, cid ClientID) {
	c.id = ClientID(askUniqueIntID("Introduzca un nuevo ID: ", func(n int64) bool { return g.ownerIDexists(ClientID(n)) }))
	for _, v := range g.vehicles {
		if cid == v.ownerID {
			v.ownerID = c.id
		}
	}
}

// modificación de cliente
func (g *Garage) modifyClient() {
	var option int
	var client *Client
	var cid ClientID

	if len(g.clients) == 0 {
		fmt.Printf("(No hay Clientes que modificar)\n\n")
		return
	}
	fmt.Printf("=== Modificación de Cliente ===\n\n")
	fmt.Printf("Lista de IDs de los clientes\n")
	for _, c := range g.clients {
		fmt.Printf("	- ID: %d. Name: %s\n", c.id, c.name)
	}
	fmt.Println()
	cid = ClientID(askUniqueIntID("Introduzca el ID del cliente a modificar: ", func(n int64) bool { return !g.ownerIDexists(ClientID(n)) }))
	client = g.getClientByID(cid)
	option = polyAskMenuInt("Seleccione lo que desee modificar.", 1)
	switch option {
	case 1:
		g.modifyClientID(client, cid)
	case 2:
		client.name = askStr("Introduzca un nuevo nombre: ")
	case 3:
		client.phone = askStr("Introduzca nuevo teléfono: ")
	case 4:
		client.email = askStr("Introduzca nuevo email: ")
	}
}

// modificar matrícula de un coche
func (g *Garage) modifIDvehicle(v *Vehicle, vid VehicleID) {
	var c *Client

	v.id = VehicleID(askUniqueStrID("Introduzca una nueva matrícula: ", func(s string) bool { return g.vidExists(VehicleID(s)) }))
	c = g.getClientByID(v.ownerID)
	for idx, idv := range c.vehicles { // actualizar matriculas asociadas a cliente
		if idv == vid {
			c.vehicles[idx] = v.id
		}
	}
	for _, issueID := range v.issues { // actualizar las incidencias asociadas a la nueva matricula
		issue := g.getIssueByID(issueID)
		issue.vehicleID = v.id
	}
	for _, s := range g.slots { // actualizar la matrícula en los slots en los que está
		if s.vehicleID != nil && *s.vehicleID == vid {
			*s.vehicleID = v.id
		}
	}
}

// modificación de vehículos
func (g *Garage) modifyVehicle() {
	var option int
	var vehicle *Vehicle
	var vid VehicleID

	if len(g.vehicles) == 0 {
		fmt.Printf("(No hay vehículos por modificar)\n\n")
		return
	}
	fmt.Printf("=== Modificación de un Vehículo ===\n\n")
	fmt.Printf("Lista de Vehículos\n")
	for _, v := range g.vehicles {
		fmt.Printf("	- %s\n", v.id)
	}
	fmt.Println()
	vid = VehicleID(askUniqueStrID("Introduzca la matrícula del vehículo a modificar: ", func(s string) bool { return !g.vidExists(VehicleID(s)) }))
	vehicle = g.getVByID(vid)
	option = polyAskMenuInt("Seleccione lo que desee modificar.", 2)
	switch option {
	case 1:
		g.modifIDvehicle(vehicle, vid)
	case 2:
		vehicle.brand = askStr("Introduzca nueva Marca: ")
	case 3:
		vehicle.model = askStr("Introduzca nuevo Modelo: ")
	}
}

// Modificar el ID de una incidencia
func (g *Garage) modifyIDissue(i *Incidence, oldid IncidenceID) {
	i.id = IncidenceID(askUniqueIntID("Introduzca un nuevo ID: ", func(n int64) bool { return g.issueIDexists(IncidenceID(n)) }))
	for _, v := range g.vehicles {
		for idx := range v.issues {
			if v.issues[idx] == oldid {
				v.issues[idx] = i.id
			}
		}
	}
	for _, m := range g.mechanics {
		for idx := range m.issues {
			if m.issues[idx] == oldid {
				m.issues[idx] = i.id
			}
		}
	}
}

// cambiar de vehiculo asociado a una incidencia
func (g *Garage) modifyAsignVIssue(i *Incidence) {
	var oldv *Vehicle
	var newv *Vehicle

	if len(g.vehicles) == 0 {
		fmt.Printf("(No hay vehículos para asignar)\n\n")
		return
	}
	oldv = g.getVByID(i.vehicleID)
	fmt.Printf("Lista de vehículos a asignar.\n")
	for _, v := range g.vehicles {
		if oldv.id != v.id {
			fmt.Printf("	- %s\n", v.id)
		}
	}
	fmt.Println()
	i.vehicleID = VehicleID(askUniqueStrID("Introduzca la mátricula que desea asignar: ", func(s string) bool { return !g.vidExists(VehicleID(s)) }))
	newv = g.getVByID(i.vehicleID)
	newv.issues = append(newv.issues, i.id)
	for idx, issueID := range oldv.issues {
		if issueID == i.id {
			oldv.issues = append(oldv.issues[:idx], oldv.issues[idx+1:]...)
		}
	}
}

// modificar el ID de un mecánico
func (g *Garage) modifIDMech(m *Mechanic, oldid MechanicID) {
	m.id = MechanicID(askUniqueIntID("Introduzca un nuevo ID: ", func(n int64) bool { return g.mechIDexists(MechanicID(n)) }))
	for i, issue := range g.issues {
		for k, mid := range issue.mechanics {
			if mid == oldid {
				g.issues[i].mechanics[k] = m.id
			}
		}
	}
}

// modificación de un mecánico
func (g *Garage) modifyMech() {
	var mech *Mechanic
	var mid MechanicID
	var option int

	if len(g.mechanics) == 0 {
		fmt.Printf("(No hay mecánicos disponibles)\n\n")
		return
	}
	fmt.Printf("=== Modificación de mecánico ===\n\n")
	fmt.Printf("Lista de mecánicos:\n")
	for _, m := range g.mechanics {
		fmt.Printf("	- ID: %d. Name: %s\n", m.id, m.name)
	}
	fmt.Println()
	mid = MechanicID(askUniqueIntID("Introduzca el ID del mecánico a modificar: ", func(n int64) bool { return !g.mechIDexists(MechanicID(n)) }))
	option = polyAskMenuInt("Seleccione lo que desee modificar.", 5)
	fmt.Println()
	mech = g.getMechByID(mid)
	switch option {
	case 1:
		g.modifIDMech(mech, mid)
	case 2:
		mech.name = askStr("Introduzca el nuevo nombre: ")
	case 3:
		mech.skill = SkillType(polyAskMenuStr("Seleccione especialidad.", 1))
	case 4:
		mech.status = MechStatus(polyAskMenuStr("Seleccione el estado.", 2))
	}

}

// cambiar mecánicos en una incidencia
func (g *Garage) asignMechs(issue *Incidence) {
	var option int

	option = polyAskMenuInt("Seleccione si desea asignar o quitar mecánicos", 4)
	fmt.Println()
	switch option {
	case 1:
		g.assignMechsToIssue(issue)
	case 2:
		g.removeMechsFromIssue(issue)
	}
}

// Modificar incidencia
func (g *Garage) modifyIssue() {
	var issueID IncidenceID
	var issue *Incidence
	var option int

	if len(g.issues) == 0 {
		fmt.Printf("(No hay incidencias aún)\n\n")
		return
	}
	fmt.Printf("=== Modificación de Incidencia ===\n\n")
	fmt.Printf("Lista de Incidencias\n")
	for _, i := range g.issues {
		fmt.Printf("	- ID: %d\n", i.id)
	}
	fmt.Println()
	issueID = IncidenceID(askUniqueIntID("Introduzca el ID de la incidencia a modificar: ", func(n int64) bool { return !g.issueIDexists(IncidenceID(n)) }))
	issue = g.getIssueByID(issueID)
	option = polyAskMenuInt("Seleccione lo que desee modificar", 3)
	fmt.Println()
	switch option {
	case 1:
		g.modifyIDissue(issue, issueID)
	case 2:
		g.modifyAsignVIssue(issue)
	case 3:
		issue.kind = IssueType(polyAskMenuStr("Seleccione tipo:", 3))
	case 4:
		issue.prio = Priority(polyAskMenuStr("Seleccione prioridad:", 4))
	case 5:
		issue.status = IssueStatus(polyAskMenuStr("Seleccione el estado:", 5))
	case 6:
		issue.description = askStr("Redacte la descripción:\n")
	case 7:
		g.asignMechs(issue)
	}
}

// meter un vehiculo a una plaza del taller
func (g *Garage) addV() {
	var vid VehicleID

	if g.availableSlots() == 0 {
		fmt.Printf("No hay plazas disponibles, espere o saque un vehículo antes.\n\n")
		return
	}
	fmt.Printf("Vehículos no asignados a una plaza:\n")
	for _, v := range g.vehicles {
		if !g.vIsInSlot(v.id) {
			fmt.Printf("	- %s\n", v.id)
		}
	}
	fmt.Println()
	vid = VehicleID(askUniqueStrID("Introduzca la matrícula del vehiculo a asignar: ", func(s string) bool { return !g.vidExists(VehicleID(s)) }))
	g.assignVtoSlot(vid)

}

// extraer un vehiculo de una plaza del taller
func (g *Garage) extractV() {
	var vid VehicleID

	if g.availableSlots() == len(g.slots) {
		fmt.Printf("No hay vehículos que sacar. Todas las plazas libres.\n\n")
		return
	}
	fmt.Printf("Vehículos situados en una plaza:\n")
	for _, v := range g.vehicles {
		if g.vIsInSlot(v.id) {
			fmt.Printf("	- %s\n", v.id)
		}
	}
	fmt.Println()
	vid = VehicleID(askUniqueStrID("Introduzca la matrícula del vehículo a extraer: ", func(s string) bool { return !g.vidExists(VehicleID(s)) }))
	g.extractVfromSlot(vid)
}

// modificar plazas el taller
func (g *Garage) modifySlots() {
	var option int

	if len(g.slots) == 0 {
		fmt.Printf("(No hay plazas disponibles, contrate mecánico)\n\n")
		return
	}
	fmt.Printf("=== Asignación de plazas del taller ===\n\n")
	for {
		option = polyAskMenuInt("Seleccione si quiere asignar o sacar un coche de una plaza:", 6)
		switch option {
		case 1:
			g.addV()
		case 2:
			g.extractV()
		case 3:
			return
		}
	}
}

func (g *Garage) assignVtoSlot(vid VehicleID) {
	var vehicle *Vehicle

	vehicle = g.getVByID(vid)
	if len(vehicle.issues) == 0 {
		return
	}

	for _, s := range g.slots {
		if s.vehicleID == nil {
			s.vehicleID = &vid
			return
		}
	}

}

// sacar un vehículo de la plaza del taller
func (g *Garage) extractVfromSlot(vid VehicleID) {
	for _, s := range g.slots {
		if s.vehicleID != nil && *s.vehicleID == vid {
			s.vehicleID = nil
		}
	}
}
