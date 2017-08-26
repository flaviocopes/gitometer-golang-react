// // eslint-disable-next-line
// import React, { Component } from 'react'

// import logo from './logo.svg'
// import './App.css'

// const App = () => (
//   <div className="App">
//     <div className="App-header">
//       <img src={logo} className="App-logo" alt="logo" />
//       <h2>Welcome to React</h2>
//     </div>
//     <p className="App-intro">
//       To gset started, edit <code>src/App.js</code> and save to reload.
//     </p>
//   </div>
// )

// export default App

import React from 'react'
import Main from './components/Main'

const App = () =>
  (
    <div>
      <div className="all-wrapper menu-side">
        <div className="layout-w">
          <div className="content-w">
            <div className="content-i">
              <div className="content-box">
                <div className="row">
                  <div className="col-lg-12">
                    <Main />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )

export default App
