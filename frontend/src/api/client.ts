import type { AuditEvent, Bot, BindingSession, Group, GroupMessage, OutboundMessage } from './types'

export const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export class Client {
  constructor(private readonly token: string) {}
  private async request<T>(path:string, init:RequestInit={}):Promise<T>{
    const response=await fetch(`${API_BASE}${path}`,{...init,headers:{'Content-Type':'application/json',Authorization:`Bearer ${this.token}`,...init.headers}})
    if(!response.ok){const body=await response.json().catch(()=>({message:response.statusText}));throw new Error(body.message ?? response.statusText)}
    if(response.status===204)return undefined as T
    return response.json() as Promise<T>
  }
  getBot(){return this.request<Bot>('/api/v1/bot')}
  startBinding(){return this.request<BindingSession>('/api/v1/bot/bindings',{method:'POST'})}
  cancelBinding(){return this.request<void>('/api/v1/bot/bindings/current',{method:'DELETE'})}
  listGroups(){return this.request<Group[]>('/api/v1/groups')}
  updateGroup(id:string,enabled:boolean){return this.request<Group>(`/api/v1/groups/${encodeURIComponent(id)}`,{method:'PATCH',body:JSON.stringify({enabled})})}
  listMessages(groupId?:string){const q=groupId?`?groupId=${encodeURIComponent(groupId)}`:'';return this.request<GroupMessage[]>(`/api/v1/messages${q}`)}
  sendMessage(groupId:string,content:string,key=crypto.randomUUID()){return this.request<OutboundMessage>(`/api/v1/groups/${encodeURIComponent(groupId)}/messages`,{method:'POST',headers:{'Idempotency-Key':key},body:JSON.stringify({content})})}
  listAudit(){return this.request<AuditEvent[]>('/api/v1/audit-events')}
  mockLogin(){return this.request('/api/v1/mock/login',{method:'POST'})}
  mockLogout(){return this.request('/api/v1/mock/logout',{method:'POST'})}
  injectMockMessage(group:Group,text:string){return this.request('/api/v1/mock/messages',{method:'POST',body:JSON.stringify({providerGroupId:group.providerGroupId,groupName:group.name,senderId:'mock-member',senderName:'群成员',text})})}
}
