// eslint-disable-next-line
import React, { Component } from 'react'
import axios from 'axios'
import { Switch, Route } from 'react-router-dom'
import RepositoriesList from './RepositoriesList'
import SingleRepository from './SingleRepository/Main'

class Repositories extends Component {
  constructor(props) {
    super(props)
    this.state = { data: null }
  }

  updateData = () => {
    axios.get('http://localhost:8000/api/repo/avelino/awesome-go').then((resp) => {
      this.setState({ data: resp.data })
      this.forceUpdate()
    })
  }

  render() {
    return (
      <Switch>
        <Route exact path="/repositories" component={RepositoriesList} />
        <Route
          path="/repositories/:name"
          render={() => <SingleRepository data={this.state.data} updateData={this.updateData} />}
        />
      </Switch>
    )
  }
}

export default Repositories
