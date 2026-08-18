export type BotStatus = 'unbound' | 'awaiting_scan' | 'connecting' | 'online' | 'reconnecting' | 'offline' | 'error'
export interface Bot { id:string; publicAccountId:string|null; displayName:string|null; status:BotStatus; statusReason:string|null; qrCode:string|null; qrExpiresAt:string|null; lastSeenAt:string|null; updatedAt:string }
export interface BindingSession { id:string; status:string; qrCode:string|null; createdAt:string; expiresAt:string; completedAt:string|null; failureReason:string|null }
export interface Group { id:string; providerGroupId:string; name:string; memberCount:number|null; enabled:boolean; available:boolean; discoveredAt:string; updatedAt:string }
export interface GroupMessage { id:string; providerMessageId:string; groupId:string; senderId:string; senderName:string; direction:'inbound'|'outbound'; messageType:'text'|'unsupported'|'system'; content:string; occurredAt:string; receivedAt:string; processingStatus:'received'|'ignored'|'sent'|'failed'; selfMessage:boolean }
export interface OutboundMessage { id:string; idempotencyKey:string; groupId:string; content:string; status:'pending'|'sending'|'succeeded'|'failed'; providerMessageId:string|null; failureReason:string|null; createdAt:string; updatedAt:string }
export interface AuditEvent { id:string; actor:string; action:string; targetType:string; targetId:string; result:'success'|'failure'; detail:string|null; createdAt:string }
export interface APIError { code:string; message:string }
