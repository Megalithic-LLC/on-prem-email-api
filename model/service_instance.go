package model

import (
	"crypto/md5"
	"time"
)

type ServiceInstance struct {
	ID              string     `json:"id" gorm:"primary_key;type:char(20)"`
	AgentID         string     `json:"agent" gorm:"type:char(20);index"`
	ServiceID       string     `json:"service" gorm:"type:char(20);index"`
	AccountIDs      []string   `json:"accounts" gorm:"-"`
	DomainIDs       []string   `json:"domains" gorm:"-"`
	PlanID          string     `json:"plan" gorm:"type:char(20);index"`
	CreatedByUserID string     `json:"createdBy" gorm:"type:char(20)"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt" gorm:"index"`
}

func (self ServiceInstance) Hash() []byte {
	hasher := md5.New()
	hasher.Write([]byte(self.ID))
	hasher.Write([]byte(self.AgentID))
	hasher.Write([]byte(self.ServiceID))
	hasher.Write([]byte(self.PlanID))
	hasher.Write([]byte(self.CreatedByUserID))

	createdAtAsBinary, _ := self.CreatedAt.MarshalBinary()
	hasher.Write(createdAtAsBinary)

	updatedAtAsBinary, _ := self.UpdatedAt.MarshalBinary()
	hasher.Write(updatedAtAsBinary)

	if self.DeletedAt != nil {
		deletedAtAsBinary, _ := self.DeletedAt.MarshalBinary()
		hasher.Write(deletedAtAsBinary)
	}

	return hasher.Sum(nil)
}
