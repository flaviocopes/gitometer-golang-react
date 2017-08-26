// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'

const Heading = props =>
  (<div className="up-head-w">
    <div className="up-main-info" style={{ textAlign: 'center' }}>
      <div className="user-avatar-w">
        <div className="user-avatar">
          <img alt="" src={props.avatar} />
        </div>
      </div>
      <br />
      <h1 className="up-header">
        <code>
          <a href={props.url}>
            {props.name}
          </a>
        </code>
      </h1>

      <br />
      <br />
      <br />
      <h3>
        {props.description}
      </h3>
      <br />
      <br />
    </div>
  </div>)

Heading.propTypes = {
  name: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
  avatar: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
}

export default Heading
