import { render, screen } from '@testing-library/react'
import { MessagePanel } from './MessagePanel'
import type { Group, GroupMessage } from '../api/types'
const group={id:'g1',providerGroupId:'p1',name:'研发群',memberCount:8,enabled:true,available:true,discoveredAt:new Date().toISOString(),updatedAt:new Date().toISOString()} satisfies Group
const message={id:'m1',providerMessageId:'pm1',groupId:'g1',senderId:'u1',senderName:'Marc',direction:'inbound',messageType:'text',content:'hello room',occurredAt:new Date().toISOString(),receivedAt:new Date().toISOString(),processingStatus:'received',selfMessage:false} satisfies GroupMessage
test('renders group message context',()=>{render(<MessagePanel messages={[message]} groups={[group]}/>);expect(screen.getByText('hello room')).toBeInTheDocument();expect(screen.getByText('研发群')).toBeInTheDocument()})
