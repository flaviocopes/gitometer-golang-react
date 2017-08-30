// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'
import Graph from './Graph'
import ValueAndPercentage from './ValueAndPercentage'

const StarsOverTime = props =>
  (<div className="up-contents">
    <div className="os-tabs-w">
      <h5 className="element-header center">
        Stars over time{' '}
        <img
          alt="&#x1f31f;"
          className="emojione"
          src="https://cdn.jsdelivr.net/emojione/assets/3.0/png/32/1f31f.png"
          title=":star2:"
        />
      </h5>
      <div className="tab-content">
        <div className="tab-pane active">
          <div className="element-content">
            <div className="row">
              <ValueAndPercentage text={'Total stars'} total={props.total} />
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
          <div className="el-chart-w">
            <Graph labels={props.graph_labels} data={props.graph_data} label={props.graph_label} />
          </div>
        </div>
        <div className="tab-pane" id="tab_conversion" />
      </div>
    </div>
  </div>)

StarsOverTime.propTypes = {
  total: PropTypes.number.isRequired,
  count_last_12_months: PropTypes.number.isRequired,
  count_last_4_weeks: PropTypes.number.isRequired,
  count_last_week: PropTypes.number.isRequired,
  graph_label: PropTypes.string,
  graph_labels: PropTypes.arrayOf(PropTypes.string),
  graph_data: PropTypes.arrayOf(PropTypes.number),
}

StarsOverTime.defaultProps = {
  graph_label: null,
  graph_labels: null,
  graph_data: null,
}

export default StarsOverTime
