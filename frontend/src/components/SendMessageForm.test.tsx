import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { SendMessageForm } from './SendMessageForm'
import type { Group } from '../api/types'
const group={id:'g1',providerGroupId:'p1',name:'研发群',memberCount:8,enabled:true,available:true,discoveredAt:new Date().toISOString(),updatedAt:new Date().toISOString()} satisfies Group
test('submits trimmed text',async()=>{const onSend=vi.fn().mockResolvedValue(undefined);render(<SendMessageForm group={group} onSend={onSend} busy={false}/>);await userEvent.type(screen.getByLabelText('消息内容'),'  hello  ');await userEvent.click(screen.getByRole('button',{name:/发送消息/}));expect(onSend).toHaveBeenCalledWith('hello')})
