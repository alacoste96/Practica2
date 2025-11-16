package main

import "time"

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
	eta       time.Duration // fecha estimada de salida
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
