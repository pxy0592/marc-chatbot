import { useEffect } from 'react'
import { API_BASE } from '../api/client'

export function useBotEvents(token:string,onEvent:()=>void){
  useEffect(()=>{
    if(!token)return
    const controller=new AbortController()
    let retry:number|undefined
    const connect=async()=>{
      try{
        const response=await fetch(`${API_BASE}/api/v1/events`,{headers:{Authorization:`Bearer ${token}`},signal:controller.signal})
        if(!response.ok||!response.body)throw new Error('event stream unavailable')
        const reader=response.body.getReader();const decoder=new TextDecoder();let buffer=''
        while(true){const {value,done}=await reader.read();if(done)break;buffer+=decoder.decode(value,{stream:true});let split=buffer.indexOf('\n\n');while(split>=0){const block=buffer.slice(0,split);buffer=buffer.slice(split+2);if(block.includes('data:'))onEvent();split=buffer.indexOf('\n\n')}}
      }catch{if(!controller.signal.aborted)retry=window.setTimeout(connect,2000)}
    }
    void connect()
    return()=>{controller.abort();if(retry)window.clearTimeout(retry)}
  },[token,onEvent])
}
