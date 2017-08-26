// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'
import { Line } from 'react-chartjs-2'

const Graph = (props) => {
  const data = {
    labels: props.labels,
    datasets: [
      {
        label: props.label,
        fill: true,
        lineTension: 0.4,
        borderColor: '#8f1cad',
        borderCapStyle: 'butt',
        borderDash: [],
        borderDashOffset: 0.0,
        borderJoinStyle: 'miter',
        pointBorderColor: '#fff',
        pointBackgroundColor: '#2a2f37',
        pointBorderWidth: 2,
        pointHoverRadius: 6,
        pointHoverBackgroundColor: '#FC2055',
        pointHoverBorderColor: '#fff',
        pointHoverBorderWidth: 2,
        pointRadius: 4,
        pointHitRadius: 5,
        data: props.data,
        spanGaps: false,
      },
    ],
  }

  return <Line data={data} />
}

Graph.propTypes = {
  label: PropTypes.string.isRequired,
  labels: PropTypes.arrayOf(PropTypes.string).isRequired,
  data: PropTypes.arrayOf(PropTypes.number).isRequired,
}

export default Graph
