import { render, screen } from '@testing-library/react'
import { BindingPanel } from './BindingPanel'
import type { Bot } from '../api/types'
const bot:Bot={id:'primary',publicAccountId:null,displayName:null,status:'awaiting_scan',statusReason:null,qrCode:'mock-code',qrExpiresAt:null,lastSeenAt:null,updatedAt:new Date().toISOString()}
test('shows binding qr and status',()=>{render(<BindingPanel bot={bot} busy={false} onBind={()=>{}} onCancel={()=>{}}/>);expect(screen.getByText('等待扫码')).toBeInTheDocument();expect(screen.getByRole('button',{name:'取消绑定'})).toBeInTheDocument()})
