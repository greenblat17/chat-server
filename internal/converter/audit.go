package converter

import "github.com/greenblat17/chat-server/internal/model"

// ToAuditFromEntity convert entity to audit model
func ToAuditFromEntity(entity model.EntityType, action string) *model.Audit {
	return &model.Audit{
		Entity: entity,
		Action: action,
	}
}
