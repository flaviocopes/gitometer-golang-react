// eslint-disable-next-line
import React, { Component } from 'react'
import PropTypes from 'prop-types'
import RepositorySummary from './RepositorySummary'
import AddNew from './AddNew'

const h1 = {
  textAlign: 'center',
  marginTop: '30px',
  fontWeight: '600',
  lineHeight: '1.1',
  color: '#334152',
  fontSize: '60px',
}

class RepositoriesList extends Component {
  componentDidMount() {
    if (!this.props.data) {
      this.props.getRepositoriesData()
    }
  }

  render() {
    if (this.props.data === null) {
      return (
        <div>
          <div className="section-heading">
            <img src={'img/logo.png'} className="center" alt="" style={{ display: 'block' }} />
            <h1 style={h1}>Gitometer</h1>
          </div>
          <div className="padded-lg">
            <div className="projects-list" style={{ textAlign: 'center' }}>
              Loading...
            </div>
          </div>
        </div>
      )
    }

    return (
      <div>
        <div className="section-heading">
          <img src={'img/logo.png'} className="center" alt="" style={{ display: 'block' }} />
          <h1 style={h1}>Gitometer</h1>
        </div>
        <div className="padded-lg">
          <div className="projects-list">
            {this.props.data.map(r =>
              (<RepositorySummary
                key={r.name}
                name={r.name}
                owner={r.ownerName}
                totalStars={r.totalStars}
                id={r.id}
              />),
            )}
            <br />
            <AddNew addNewRepository={this.props.addNewRepository} />
          </div>
        </div>
      </div>
    )
  }
}

RepositoriesList.propTypes = {
  addNewRepository: PropTypes.func.isRequired,
  getRepositoriesData: PropTypes.func.isRequired,
  data: PropTypes.arrayOf(
    PropTypes.shape({
      name: PropTypes.string,
      totalStars: PropTypes.number,
      id: PropTypes.number,
    }),
  ),
}

RepositoriesList.defaultProps = {
  data: null,
}

export default RepositoriesList
