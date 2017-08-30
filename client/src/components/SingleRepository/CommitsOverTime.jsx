// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'
import ValueAndPercentage from './ValueAndPercentage'

const CommitsOverTime = props =>
  (<div className="up-contents">
    <div className="os-tabs-w">
      <h5 className="element-header center">Commits on the default branch</h5>
      <div className="tab-content">
        <div className="tab-pane active">
          <div className="element-content">
            <div className="row">
              <ValueAndPercentage text={'Total commits'} total={props.total} />
              <ValueAndPercentage
                text={'Last 12 months'}
                total={props.total}
                partial={props.count_last_12_months}
              />
              <ValueAndPercentage
                text={'Last 4 weeks'}
                total={props.total}
                partial={props.count_last_4_weeks}
              />
              <ValueAndPercentage
                text={'Last week'}
                total={props.total}
                partial={props.count_last_week}
              />
            </div>
          </div>
        </div>
        <div className="tab-pane" id="tab_conversion" />
      </div>
    </div>
  </div>)

CommitsOverTime.propTypes = {
  total: PropTypes.number.isRequired,
  count_last_12_months: PropTypes.number.isRequired,
  count_last_4_weeks: PropTypes.number.isRequired,
  count_last_week: PropTypes.number.isRequired,
}

export default CommitsOverTime
