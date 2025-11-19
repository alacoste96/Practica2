package main

import (
	"fmt"
	"time"
)

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
	go g.assignVtoSlot(vid)
}

// modificar plazas el taller
func (g *Garage) modifySlots() {

	if len(g.slots) == 0 {
		fmt.Printf("(No hay plazas disponibles, contrate mecánico)\n\n")
		return
	}
	g.addV()
}

// devuelve un mecánico que está libre
func (g *Garage) getFreeMech() *Mechanic {
	for _, m := range g.mechanics {
		if isFree(m) {
			return m
		}
	}
	return nil
}

func (g *Garage) autoAsignMtoIssues(m *Mechanic, issues []*Incidence) {

	for _, i := range issues {
		m.issues = append(m.issues, i.id)
		i.mechanics = append(i.mechanics, m.id)
	}
}

func (g *Garage) waitForFreeMech() *Mechanic {
	var m *Mechanic
	for {
		m = g.getFreeMech()
		if m != nil {
			return m
		}
		time.Sleep(1 * time.Second)
	}
}

func setPrioAndStatus(issues []*Incidence) {
	for _, i := range issues {
		i.prio = HIGH
		i.status = INPROCESS
	}
}

func closeIssues(issues []*Incidence) {
	for _, i := range issues {
		i.status = CLOSED
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

func (g *Garage) assignVtoSlot(vid VehicleID) {
	vehicle := g.getVByID(vid)
	if vehicle == nil || len(vehicle.issues) == 0 {
		return
	}

	// 1) Obtener mecánico principal (esperando si hace falta)
	mainMech := g.obtainMechanicForVehicle(vid)

	// 2) Meter vehículo en una plaza
	if !g.placeVehicleInSlot(vid) {
		// si no hubiera slot libre (por seguridad)
		return
	}

	// 3) Recoger incidencias y calcular ETA
	issues, eta := g.collectIssuesAndEta(vehicle)
	vehicle.eta = eta

	// 4) Asignar mecánico principal
	g.autoAsignMtoIssues(mainMech, issues)

	// 5) Posible segundo mecánico
	var extraMech *Mechanic
	if eta > 15*time.Second {
		extraMech = g.obtainSecondMechanic(vid)
		setPrioAndStatus(issues)
		g.autoAsignMtoIssues(extraMech, issues)
	}

	// 6) Simular reparación
	time.Sleep(eta)

	// 7) Cerrar incidencias y liberar recursos
	closeIssues(issues)
	g.extractVfromSlot(vid)

	for _, inc := range issues {
		delIssueFromMech(mainMech, inc)
		if extraMech != nil {
			delIssueFromMech(extraMech, inc)
		}
		delIssueFromV(vehicle, inc)
	}
}

// 1) Mecánico principal (usa vpool + wait)
func (g *Garage) obtainMechanicForVehicle(vid VehicleID) *Mechanic {
	if m := g.getFreeMech(); m != nil {
		return m
	}

	// nadie libre -> coche a la cola de espera
	g.vpool = append(g.vpool, vid)
	m := g.waitForFreeMech()
	g.vpool = removeVID(g.vpool, vid)
	return m
}

// 2) Colocar vehículo en un slot libre
func (g *Garage) placeVehicleInSlot(vid VehicleID) bool {
	for i := range g.slots {
		if g.slots[i].vehicleID == nil {
			g.slots[i].vehicleID = &vid
			return true
		}
	}
	return false
}

// 3) Construir slice de incidencias y calcular ETA
func (g *Garage) collectIssuesAndEta(v *Vehicle) ([]*Incidence, time.Duration) {
	var (
		issues []*Incidence
		eta    time.Duration
	)

	for _, issueID := range v.issues {
		inc := g.getIssueByID(issueID)
		if inc == nil {
			continue
		}
		issues = append(issues, inc)

		switch inc.kind {
		case MECHTYPE:
			eta += 5 * time.Second
		case ELECTRICTYPE:
			eta += 7 * time.Second
		default:
			eta += 11 * time.Second
		}
	}
	return issues, eta
}

// 4) Segundo mecánico (libre o contratado)
func (g *Garage) obtainSecondMechanic(vid VehicleID) *Mechanic {
	if m := g.getFreeMech(); m != nil {
		return m
	}
	req := HireRequest{
		vid:   vid,
		Reply: make(chan *Mechanic),
	}
	g.hirereqs <- req
	return <-req.Reply
}
