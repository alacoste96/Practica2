package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

/*
------------- To do ------------------
- Menús


*/
// -------------------------------
// -------- MIS CONSTANTES -------
// -------------------------------

type IssueType string

const (
	MECHTYPE     IssueType = "Mecánica"
	ELECTRICTYPE IssueType = "Eléctrica"
	BODYTYPE     IssueType = "Carrocería"
)

type Priority string

const (
	LOW    Priority = "Baja"
	MEDIUM Priority = "Media"
	HIGH   Priority = "Alta"
)

type IssueStatus string

const (
	OPEN      IssueStatus = "Abierta"
	INPROCESS IssueStatus = "En proceso"
	CLOSED    IssueStatus = "Cerrada"
)

type SkillType string

const (
	MECHSKILL     SkillType = "Mecánica"
	ELECTRICSKILL SkillType = "Eléctrica"
	BODYSKILL     SkillType = "Carrocería"
)

type MechStatus string

const (
	ACTIVE   MechStatus = "Activo"
	INACTIVE MechStatus = "De baja"
)

// enumero tipos de datos para entenderlos mejor como humano
type ClientID int64
type MechanicID int64
type IncidenceID int64
type VehicleID string

// --------------------------------------
// ---------- MIS ESTRUCTURAS -----------
// --------------------------------------

// Client representa a un cliente
type Client struct {
	id       ClientID
	name     string
	phone    string
	email    string
	vehicles []VehicleID
	// los vehículos del cliente se pueden consultar por índice
}

// Mechanic representa un mecánico.
type Mechanic struct {
	id         MechanicID
	name       string
	skill      SkillType
	experience int
	status     MechStatus
	issues     []IncidenceID
}

// Vehicle representa un vehículo en el taller.
// OwnerID mantiene la relación 1..N con Client.
type Vehicle struct {
	id        VehicleID // matrícula
	brand     string
	model     string
	ownerID   ClientID
	checkInAt time.Time
	eta       time.Time // fecha estimada de salida
	issues    []IncidenceID
}

// Incidence representa una incidencia/reparación.
// Asignación N..M con mecánicos mediante slice de IDs.
type Incidence struct {
	id          IncidenceID
	vehicleID   VehicleID
	mechanics   []MechanicID
	kind        IssueType
	prio        Priority
	description string
	status      IssueStatus
}

// Slot representa una plaza física del taller. Si VehicleID == nil, está libre.
type Slot struct {
	number    int
	vehicleID *VehicleID
}

// mi taller, de él recaerá el peso de muchos métodos para su propia gestión
// La capacidad debe seguir la regla: 2 plazas por mecánico activo.
type Garage struct {
	clients   []*Client
	vehicles  []*Vehicle
	mechanics []*Mechanic
	issues    []*Incidence
	slots     []*Slot
}

// ---------------------------------------------------
// ------------ DISPLAYS DE VISUALIZACIÓN ------------
// ---------------------------------------------------

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

// --------------------------------------------
// ----------- GESTIÓN DEL TALLER -------------
//---------------------------------------------

// -------------------------------
// ---- CREACIÓN DE ENTIDADES ----
// -------------------------------

// Constructor default de un nuevo taller
func newGarage() *Garage {
	return &Garage{
		clients:   nil,
		vehicles:  nil,
		mechanics: nil,
		issues:    nil,
		slots:     nil}
}

// CREAR un CLIENTE
func (g *Garage) newClient() {
	var c Client

	c.id = ClientID(askUniqueIntID("ID cliente (>0): ", func(n int64) bool { return g.ownerIDexists(ClientID(n)) }))
	fmt.Printf("Nombre: ")
	fmt.Scanf("%s", &c.name)
	fmt.Printf("Teléfono: ")
	fmt.Scanf("%s", &c.phone)
	fmt.Printf("Email: ")
	fmt.Scanf("%s", &c.email)
	g.clients = append(g.clients, &c)
}

// CREAR VEHÍCULO
func (g *Garage) newVehicle() {
	var v Vehicle
	var client *Client

	if len(g.clients) == 0 {
		fmt.Printf("No hay aún clientes registrados, registre uno primero.\n")
		return
	}
	fmt.Printf("¿Cuál es el dueño? escoja un ID:\n\n")
	for _, c := range g.clients {
		fmt.Printf("	- ID: %d. Cliente: %s\n", c.id, c.name)
	}
	v.ownerID = ClientID(askUniqueIntID("Dueño: ", func(n int64) bool { return !g.ownerIDexists(ClientID(n)) }))
	v.id = VehicleID(askUniqueStrID("Matrícula: ", func(s string) bool { return g.vidExists(VehicleID(s)) }))
	client = g.getClientByID(v.ownerID)
	client.vehicles = append(client.vehicles, v.id)
	fmt.Printf("Marca: ")
	fmt.Scanf("%s", &v.brand)
	fmt.Printf("Modelo: ")
	fmt.Scanf("%s", &v.model)
	v.checkInAt = time.Now().Add(time.Duration(-30) * time.Hour)
	v.eta = time.Now().AddDate(0, 0, 3)
	g.vehicles = append(g.vehicles, &v)
}

// CREAR MECÁNICO
func (g *Garage) newMech() {
	var m Mechanic

	// ampliar plazas del taller
	for i := 0; i < 2; i++ {
		g.slots = append(g.slots, &Slot{
			number:    len(g.slots) + 1,
			vehicleID: nil, // puntero al ID existente en el slice (dirección estable)
		})
	}
	m.id = MechanicID(askUniqueIntID("ID del mecánico: ", func(n int64) bool { return g.mechIDexists(MechanicID(n)) }))
	fmt.Printf("Nombre: ")
	fmt.Scanf("%s", &m.name)
	fmt.Printf("Experiencia(años): ")
	fmt.Scanf("%d", &m.experience)
	m.skill = SkillType(polyAskMenuStr("Seleccione una especialidad:", 1))
	m.status = MechStatus(polyAskMenuStr("Seleccione un estado:", 2))
	g.mechanics = append(g.mechanics, &m)
}

// CREAR INCIDENCIA
func (g *Garage) newIssue() {
	var i Incidence
	var v *Vehicle

	i.id = IncidenceID(askUniqueIntID("ID de la incidencia: ", func(n int64) bool { return g.issueIDexists(IncidenceID(n)) }))
	fmt.Printf("Escriba la matrícula del coche que sufre la incidencia:\n")
	for _, v := range g.vehicles {
		fmt.Printf("	- %s\n", v.id)
	}
	fmt.Println()
	i.vehicleID = VehicleID(askUniqueStrID("Matrícula: ", func(s string) bool { return !g.vidExists(VehicleID(s)) }))
	fmt.Println()
	v = g.getVByID(i.vehicleID)
	v.issues = append(v.issues, i.id)
	i.kind = IssueType(polyAskMenuStr("Tipo de incidencia:", 3))
	fmt.Println()
	i.prio = Priority(polyAskMenuStr("Prioridad:", 4))
	fmt.Println()
	i.status = IssueStatus(polyAskMenuStr("Estado:", 5))
	fmt.Println()
	g.assignMechsToIssue(&i)
	fmt.Println()
	fmt.Printf("Redacte la Descripción de la incidencia:\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	i.description = scanner.Text()

	g.issues = append(g.issues, &i)
}

// --------------------------------------------
// --------ELIMINACION DE ENTIDADES------------
// --------------------------------------------

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

// --------------------------------------------
// ------------- MODIFICACIÓN -----------------
// --------------------------------------------

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
		fmt.Printf("No hay plazas disponibles, saque un vehículo antes.\n\n")
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

// --------------------------------------------
// ----------------- MENÚS --------------------
// --------------------------------------------

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

// --------------------------------------------
// ----------- FUNCIONES AUXILIARES -----------
// --------------------------------------------

// asignar vehículo a la primera plaza libre
func (g *Garage) assignVtoSlot(vid VehicleID) {
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
