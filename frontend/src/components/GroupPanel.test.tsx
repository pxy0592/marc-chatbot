import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { GroupPanel } from './GroupPanel'
import type { Group } from '../api/types'
const group:Group={id:'g1',providerGroupId:'p1',name:'研发群',memberCount:8,enabled:false,available:true,discoveredAt:new Date().toISOString(),updatedAt:new Date().toISOString()}
test('toggles a group',async()=>{const onToggle=vi.fn();render(<GroupPanel groups={[group]} selectedId="" onSelect={()=>{}} onToggle={onToggle}/>);await userEvent.click(screen.getByRole('checkbox',{name:'启用 研发群'}));expect(onToggle).toHaveBeenCalledWith(group)})
