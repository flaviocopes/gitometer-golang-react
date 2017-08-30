// eslint-disable-next-line
import React, { Component } from 'react'
import { Switch, Route } from 'react-router-dom'
import axios from 'axios'

import SingleRepository from './SingleRepository/Main'
import RepositoriesList from './RepositoriesList/Main'
import Settings from './Settings'

class Main extends Component {
  constructor(props) {
    super(props)
    this.state = { repositoriesData: null, repositoryData: null }
  }

  getRepositoriesData = () => {
    axios.get('http://localhost:8000/api/index').then((resp) => {
      this.setState({ repositoriesData: resp.data.repositories })
      this.forceUpdate()
    })
  }

  getRepositoryData = (owner, name) => {
    axios.get(`http://localhost:8000/api/repo/${owner}/${name}`).then((resp) => {
      let newRepositoryData = Object.assign({}, this.state.repositoryData)
      if (!newRepositoryData) {
        newRepositoryData = {}
      }
      if (!newRepositoryData[owner]) {
        newRepositoryData[owner] = {}
      }
      newRepositoryData[owner][name] = resp.data
      this.setState({ repositoryData: newRepositoryData })
    })
  }

  addNewRepository = (owner, name) => {
    axios.post('http://localhost:8000/api/repo', { owner, name }).then((resp) => {
      console.log(resp)
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
              (<RepositoriesList
                getRepositoriesData={this.getRepositoriesData}
                addNewRepository={this.addNewRepository}
                data={this.state.repositoriesData}
              />)}
          />
          <Route
            path="/repositories/:owner/:name"
            render={(props) => {
              let data = null
              if (
                this.state.repositoryData &&
                this.state.repositoryData[props.match.params.owner] &&
                this.state.repositoryData[props.match.params.owner][props.match.params.name]
              ) {
                data = this.state.repositoryData[props.match.params.owner][props.match.params.name]
              }
              return (
                <SingleRepository
                  match={props.match}
                  data={data}
                  getRepositoryData={this.getRepositoryData}
                />
              )
            }}
          />
          <Route path="/settings" component={Settings} />
        </Switch>
      </main>
    )
  }
}

export default Main
