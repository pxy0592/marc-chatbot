import { useCallback, useEffect, useMemo, useState } from 'react'
import { Client } from './api/client'
import type { AuditEvent, Bot, Group, GroupMessage } from './api/types'
import { useBotEvents } from './hooks/useBotEvents'
import { AuthGate } from './components/AuthGate'
import { BindingPanel } from './components/BindingPanel'
import { GroupPanel } from './components/GroupPanel'
import { MessagePanel } from './components/MessagePanel'
import { SendMessageForm } from './components/SendMessageForm'
import { AuditPanel } from './components/AuditPanel'

export default function App(){
  const[token,setToken]=useState(()=>sessionStorage.getItem('admin-token')??'');const[bot,setBot]=useState<Bot|null>(null);const[groups,setGroups]=useState<Group[]>([]);const[messages,setMessages]=useState<GroupMessage[]>([]);const[audit,setAudit]=useState<AuditEvent[]>([]);const[selected,setSelected]=useState('');const[busy,setBusy]=useState(false);const[error,setError]=useState('')
  const client=useMemo(()=>new Client(token),[token]);const refresh=useCallback(async()=>{if(!token)return;try{const[b,g,m,a]=await Promise.all([client.getBot(),client.listGroups(),client.listMessages(),client.listAudit()]);setBot(b);setGroups(g);setMessages(m);setAudit(a);setSelected(current=>current||g[0]?.id||'');setError('')}catch(e){setError(e instanceof Error?e.message:'加载失败')}},[client,token]);useEffect(()=>{void refresh();const timer=window.setInterval(refresh,10000);return()=>window.clearInterval(timer)},[refresh]);useBotEvents(token,refresh)
  if(!token)return <AuthGate onLogin={value=>{sessionStorage.setItem('admin-token',value);setToken(value)}}/>
  const run=async(action:()=>Promise<unknown>)=>{setBusy(true);setError('');try{await action();await new Promise(r=>setTimeout(r,100));await refresh()}catch(e){setError(e instanceof Error?e.message:'操作失败')}finally{setBusy(false)}};const selectedGroup=groups.find(g=>g.id===selected);const mock=import.meta.env.VITE_ENABLE_MOCK_CONTROLS!=='false'
  return <main className="app-shell"><header className="topbar"><div className="brand"><span>M</span><div><b>MARC CHATBOT</b><small>WECHAT CONTROL PLANE</small></div></div><div className="top-actions"><span className="connection-chip"><i/>SYSTEM ACTIVE</span><button className="icon-button" title="刷新" onClick={()=>void refresh()}>↻</button><button className="logout" onClick={()=>{sessionStorage.removeItem('admin-token');setToken('')}}>退出</button></div></header>{error&&<div className="error-banner">{error}<button onClick={()=>setError('')}>×</button></div>}<section className="dashboard"><div className="column primary-column"><BindingPanel bot={bot} busy={busy} onBind={()=>void run(()=>client.startBinding())} onCancel={()=>void run(()=>client.cancelBinding())} showMock={mock} onMockLogin={()=>void run(()=>client.mockLogin())} onMockLogout={()=>void run(()=>client.mockLogout())}/><MessagePanel messages={messages} groups={groups}/></div><div className="column side-column"><GroupPanel groups={groups} selectedId={selected} onSelect={setSelected} onToggle={group=>void run(()=>client.updateGroup(group.id,!group.enabled))}/><SendMessageForm group={selectedGroup} busy={busy} onSend={text=>run(()=>client.sendMessage(selected,text)).then(()=>undefined)}/>{mock&&selectedGroup&&<button className="inject-button" onClick={()=>void run(()=>client.injectMockMessage(selectedGroup,'这是一条模拟群消息'))}>注入一条模拟群消息</button>}<AuditPanel events={audit}/></div></section></main>
}
