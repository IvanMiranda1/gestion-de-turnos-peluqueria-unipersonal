package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type PreferenciaHoraria int

const (
	Mañana = iota
	Tarde
	Noche
)

func (s PreferenciaHoraria) String() string {
	return [...]string{"Mañana", "Tarde", "Noche"}[s]
}

func ParsePreferenciaHoraria(s string) (PreferenciaHoraria, error) {
	switch s {
	case "Mañana":
		return Mañana, nil
	case "Tarde":
		return Tarde, nil
	case "Noche":
		return Noche, nil
	default:
		return -1, fmt.Errorf("preferencia horaria no valida: %s", s)
	}
}

func IsValidPreferenciaHoraria(p PreferenciaHoraria) bool {
	switch p {
	case Mañana, Tarde, Noche:
		return true
	default:
		return false
	}
}

type TimeOfDay struct {
	Hour   int
	Minute int
}

func NewTimeOfDay(hour, minute int) (TimeOfDay, error) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return TimeOfDay{}, fmt.Errorf("hora inválida: %02d:%02d", hour, minute)
	}
	return TimeOfDay{Hour: hour, Minute: minute}, nil
}

func (t TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

// comparar si una hora es posterior a otra
func (t TimeOfDay) IsAfter(other TimeOfDay) bool {
	if t.Hour != other.Hour {
		return t.Hour > other.Hour
	}
	return t.Minute > other.Minute

}

func (t TimeOfDay) IsValid() bool {
	if t.Hour < 0 || t.Hour > 23 {
		return false
	}
	if t.Minute < 0 || t.Minute > 59 {
		return false
	}
	return true
}

func ParseTimeOfDay(input string) (TimeOfDay, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return TimeOfDay{}, fmt.Errorf("formato inválido, se esperaba 'HH:MM'")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return TimeOfDay{}, fmt.Errorf("hora inválida: %w", err)
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return TimeOfDay{}, fmt.Errorf("minuto inválido: %w", err)
	}

	return NewTimeOfDay(hour, minute)
}

func (t TimeOfDay) Equals(other TimeOfDay) bool {
	return t.Hour == other.Hour && t.Minute == other.Minute
}
