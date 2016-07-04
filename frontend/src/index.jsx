import React from 'react'
import {render} from 'react-dom'
import SockJS from 'sockjs-client'

var socket = new SockJS('/ws')

socket.onmessage = (e) => {
  console.log('message', e.data)
}

class CmdDebug extends React.Component {

  render () {
    return (
      <div>
        <textarea style={{width: '480px', height: '240px'}} ref={(e) => { this.cmd = e }} />
        <br />
        <button type='submit' onClick={(e) => {
          console.log(this.cmd.value)
          socket.send(this.cmd.value)
        }}>
          Run!
        </button>
      </div>
    )
  }
}

render(
  <CmdDebug />,
  document.getElementById('app')
)
