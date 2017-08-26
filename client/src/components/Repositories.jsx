import React from 'react'
import { Switch, Route } from 'react-router-dom'
import RepositoriesList from './RepositoriesList'
import SingleRepository from './SingleRepository/Main'

// The Roster component matches one of two different routes
// depending on the full pathname
const Repositories = () =>
  (<Switch>
    <Route exact path="/repositories" component={RepositoriesList} />
    <Route path="/repositories/:name" component={SingleRepository} />
  </Switch>)

export default Repositories
