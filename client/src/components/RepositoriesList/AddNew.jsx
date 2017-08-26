// eslint-disable-next-line
import React, { Component, ReactDOM } from 'react'
import PropTypes from 'prop-types'

class RepositorySummary extends Component {
  addNewRepository = (e) => {
    e.preventDefault()
    const owner = this.owner.value
    const name = this.name.value
    this.props.addNewRepository(owner, name)
  }

  render() {
    return (
      <form onSubmit={this.addNewRepository} style={{ textAlign: 'center', margin: '0 auto' }}>
        <br />
        <br />
        <h3>Add a new repo</h3>
        <br />
        <div className="form-group">
          <label htmlFor="owner">
            <input
              id="owner"
              ref={(c) => {
                this.owner = c
              }}
              className="form-control"
              placeholder="Enter a repo owner"
              type="text"
            />
          </label>
        </div>
        <div className="form-group">
          <label htmlFor="name">
            <input
              id="name"
              ref={(c) => {
                this.name = c
              }}
              className="form-control"
              placeholder="Enter a repo name"
              type="text"
            />
          </label>
        </div>
        <div className="buttons-w">
          <button className="btn btn-primary">Add</button>
        </div>
      </form>
    )
  }
}

RepositorySummary.propTypes = {
  addNewRepository: PropTypes.func.isRequired,
}

export default RepositorySummary
