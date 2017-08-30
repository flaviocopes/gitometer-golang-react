// eslint-disable-next-line
import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router-dom'
import Heading from './Heading'
import SummaryNumbers from './SummaryNumbers'
import CommitsOverTime from './CommitsOverTime'
import StarsOverTime from './StarsOverTime'

class SingleRepository extends Component {
  constructor(props) {
    super(props)
    this.owner = props.match.params.owner
    this.name = props.match.params.name
  }

  componentDidMount() {
    if (!this.props.data) {
      this.props.getRepositoryData(this.owner, this.name)
    }
  }

  render() {
    const data = this.props.data

    if (!data) {
      return (
        <div>
          <Link to="/">Home</Link> Sorry, but the repo was not found
        </div>
      )
    }

    const url = 'https://github.com/{data.repository.owner}/{data.repository.name}'
    let starsPerMonth = {
      labels: null,
      data_summed: null,
    }
    if (data.repository.stars_per_month) {
      starsPerMonth = JSON.parse(data.repository.stars_per_month)
    }

    return (
      <div>
        <div className="element-wrapper wrapper-dashboard">
          <div className="user-profile">
            <br />
            <Link to="/" style={{ textAlign: 'center', display: 'block' }}>Home</Link>
            <Heading
              name={data.repository.name}
              description={data.repository.description}
              url={url}
            />
            <SummaryNumbers
              since={data.repository.repository_created_months_ago}
              stars={data.repository.total_stars}
              commits={data.repository.total_commits}
            />
            <CommitsOverTime
              total={data.repository.total_commits}
              count_last_12_months={data.repository.commits_count_last_12_months}
              count_last_4_weeks={data.repository.commits_count_last_4_weeks}
              count_last_week={data.repository.commits_count_last_week}
            />
            <StarsOverTime
              total={data.repository.total_stars}
              count_last_12_months={data.repository.stars_count_last_12_months}
              count_last_4_weeks={data.repository.stars_count_last_4_weeks}
              count_last_week={data.repository.stars_count_last_week}
              graph_label={'Stars up to now'}
              graph_labels={starsPerMonth.labels}
              graph_data={starsPerMonth.data}
            />
          </div>
        </div>
      </div>
    )
  }
}

SingleRepository.propTypes = {
  getRepositoryData: PropTypes.func.isRequired,
  data: PropTypes.shape({
    repository: PropTypes.shape({
      name: PropTypes.string,
      description: PropTypes.string,
      avatar_url: PropTypes.string,
      url: PropTypes.string,
      repository_created_months_ago: PropTypes.number,
      total_stars: PropTypes.number,
      total_commits: PropTypes.number,
      commits_count_last_12_months: PropTypes.number,
      commits_count_last_4_weeks: PropTypes.number,
      commits_count_last_week: PropTypes.number,
      stars_count_last_12_months: PropTypes.number,
      stars_count_last_4_weeks: PropTypes.number,
      stars_count_last_week: PropTypes.number,
      stars_per_month: PropTypes.string,
    }),
  }),
}

SingleRepository.defaultProps = {
  data: null,
}

export default SingleRepository
