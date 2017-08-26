// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'

const SummaryNumbers = (props) => {
  const years = parseInt(props.since / 12, 10)
  const months = props.since % 12

  let since = 'This repo is on GitHub since'

  if (years > 0) {
    since += ` ${years} years`
    if (months > 0) {
      since += ` and ${months} months`
    }
  } else {
    switch (months) {
      case 0:
        since += ' less than a month'
        break
      case 1:
        since += ' 1 month'
        break
      default:
        since += ` ${months} months`
    }
  }

  return (
    <div className="up-contents center">
      <h5 className="element-header">
        {since}
      </h5>
      <div className="element-content">
        <div className="row">
          <div className="col-sm-6 b-r">
            <div className="el-tablo centered padded">
              <div className="value">
                {props.stars}
              </div>
              <div className="label">
                <img
                  alt="&#x1f31f;"
                  className="emojione"
                  src="https://cdn.jsdelivr.net/emojione/assets/3.0/png/32/1f31f.png"
                  title=":star2:"
                />{' '}
                Stars
              </div>
            </div>
          </div>
          <div className="col-sm-6">
            <div className="el-tablo centered padded">
              <div className="value">
                {props.commits}
              </div>
              <div className="label">
                <img
                  alt="&#x1f4d1;"
                  className="emojione"
                  src="https://cdn.jsdelivr.net/emojione/assets/3.0/png/32/1f4d1.png"
                  title=":bookmark_tabs:"
                />{' '}
                Commits
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

SummaryNumbers.propTypes = {
  commits: PropTypes.number.isRequired,
  stars: PropTypes.number.isRequired,
  since: PropTypes.number.isRequired,
}

export default SummaryNumbers
