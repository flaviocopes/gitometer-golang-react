// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'

const ValueAndPercentage = (props) => {
  const percentage = parseInt((props.partial / props.total) * 100, 10)

  if (props.partial === null) {
    return (
      <div className="col-sm-3">
        <div className="el-tablo padded-top-and-bottom">
          <div className="label">
            {props.text}
          </div>
          <div className="value">
            {props.total}
          </div>
        </div>
      </div>
    )
  }

  if (percentage === 0) {
    return (
      <div className="col-sm-3">
        <div className="el-tablo padded-top-and-bottom">
          <div className="label">
            {props.text}
          </div>
          <div className="value">
            {props.partial}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="col-sm-3">
      <div className="el-tablo padded-top-and-bottom">
        <div className="label">
          {props.text}
        </div>
        <div className="value">
          {props.partial}
        </div>
        <div className="trending trending-up">
          <span>
            {percentage}%
          </span>
          <i className="os-icon" />
        </div>
      </div>
    </div>
  )
}

ValueAndPercentage.propTypes = {
  text: PropTypes.string.isRequired,
  total: PropTypes.number.isRequired,
  partial: PropTypes.number,
}

ValueAndPercentage.defaultProps = {
  partial: null,
}

export default ValueAndPercentage
