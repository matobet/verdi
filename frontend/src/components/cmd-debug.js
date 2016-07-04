import React from 'react'

export default class CmdDebug extends React.Component {
  render () {
    const { ws } = this.props
    return (
      <div>
        <label htmlFor="json">Raw JSON command:</label>
        <br />
        <textarea id="json" style={{width: '480px', height: '240px'}} ref={(e) => { this.cmd = e }} />
        <br />
        <button type='submit' onClick={(e) => {
          console.log(this.cmd.value)
          ws.send(this.cmd.value)
        }}>
          Run!
        </button>
      </div>
    )
  }
}
