package main

import "time"

// Guest represents a single extended-stay hotel guest.
// This struct is used as:
// 1) The in-memory domain model
// 2) The JSON shape exposed by the API
// Mirrors the database schema to keep mapping logic consistent
type Guest struct {
	// ID is the db primary key.
	ID int `json:"id"`

	// Name is the guest’s full name as stored in the system.
	Name string `json:"name"`

	// Room assigned to the guest stored as a string to support non-numeric formats
	RoomNumber string `json:"room_number"`

	// DailyRate of guests' room stored as an integer to avoid floating-point precision issues
	DailyRate int `json:"daily_rate"`

	// CheckInDate is the date the guest started their stay
	// Used for billing period calculations
	CheckInDate time.Time `json:"check_in_date"`

	// Contact of the guest (email or phone)
	Contact string `json:"contact"`

	// WeeksStayed is a derived field, computed at runtime
	// based on CheckInDate and the current date  used for billing and notification logic
	WeeksStayed int `json:"weeks_stayed"`
}
