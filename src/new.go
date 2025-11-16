package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

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
	v.eta = 0 * time.Second
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
