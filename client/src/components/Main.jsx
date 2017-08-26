// eslint-disable-next-line
import React, { Component } from 'react'
import { Switch, Route } from 'react-router-dom'
import axios from 'axios'

import Home from './Home'
import SingleRepository from './SingleRepository/Main'
import Schedule from './Schedule'

class Main extends Component {
  constructor(props) {
    super(props)
    this.state = { data: null }
  }

  getRepositoriesData = () => {
    axios.get('http://localhost:8000/api/index').then((resp) => {
      this.setState({ data: resp.data.repositories })
      this.forceUpdate()
    })
  }

  render() {
    return (
      <main>
        <Switch>
          <Route
            exact
            path="/"
            render={() =>
              <Home getRepositoriesData={this.getRepositoriesData} data={this.state.data} />}
          />
          <Route
            path="/repositories/:name"
            render={() => <SingleRepository data={this.state.data} updateData={this.updateData} />}
          />
          <Route path="/settings" component={Settings} />
        </Switch>
      </main>
    )
  }
}

export default Main
