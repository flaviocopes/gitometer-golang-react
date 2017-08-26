// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'

import { Link } from 'react-router-dom'

const RepositorySummary = props =>
  (<div className="project-box">
    <Link to={`/repositories/${props.owner}/${props.name}`}>
      <div className="project-head">
        <div className="project-title">
          <h5 className="home-entry-stars">
            {props.totalStars}{' '}
            <img
              className="emojione"
              alt="&#x1f31f;"
              title=":star2:"
              src="https://cdn.jsdelivr.net/emojione/assets/3.0/png/32/1f31f.png"
            />
          </h5>
          <h5 className="home-entry-title">
            {props.name}
          </h5>
        </div>
      </div>
    </Link>
  </div>)

RepositorySummary.propTypes = {
  // id: PropTypes.number.isRequired,
  totalStars: PropTypes.number.isRequired,
  name: PropTypes.string.isRequired,
  owner: PropTypes.string.isRequired,
}

export default RepositorySummary
