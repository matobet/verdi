import React from 'react'
import { render } from 'react-dom'
import SockJS from 'sockjs-client'

import CmdDebug from './components/cmd-debug'

var socket = new SockJS('/ws')

socket.onmessage = (e) => {
  console.log('message', e.data)
}

render(
  <CmdDebug ws={socket} />,
  document.getElementById('app')
)
