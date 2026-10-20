package domain

type Member struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name  string `json:"name"`
	Phone string `json:"phone"`

	TimeSlots []TimeSlot `gorm:"many2many:member_time_slots;" json:"time_slots,omitempty"`
}
