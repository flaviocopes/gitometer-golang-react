// eslint-disable-next-line
import React from 'react'
import PropTypes from 'prop-types'
// import { Link } from 'react-router-dom'
import API from '../../api'
import Heading from './Heading'
import SummaryNumbers from './SummaryNumbers'
import CommitsOverTime from './CommitsOverTime'
import StarsOverTime from './StarsOverTime'

const SingleRepository = (props) => {
  const data = API.get(props.match.params.name)
  if (!data) {
    return <div>Sorry, but the repo was not found</div>
  }
  return (
    <div>
      <div className="element-wrapper wrapper-dashboard">
        <div className="user-profile">
          <Heading
            name={data.repository.repository_name}
            description={data.repository.description}
            avatar={data.organization.avatar_url}
            url={data.repository.url}
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
            graph_label={'Commits up to now'}
            graph_labels={data.monthly_data.commits_per_month.labels}
            graph_data={data.monthly_data.commits_per_month.data_summed}
          />
          <StarsOverTime
            total={data.repository.total_stars}
            count_last_12_months={data.repository.stars_count_last_12_months}
            count_last_4_weeks={data.repository.stars_count_last_4_weeks}
            count_last_week={data.repository.stars_count_last_week}
            graph_label={'Stars up to now'}
            graph_labels={data.monthly_data.stars_per_month.labels}
            graph_data={data.monthly_data.stars_per_month.data_summed}
          />
        </div>
      </div>
    </div>
  )
}

// SingleRepository.propTypes = {
//   id: PropTypes.number.isRequired,
//   stars: PropTypes.number.isRequired,
//   name: PropTypes.string.isRequired,
//   description: PropTypes.string.isRequired,
// }

export default SingleRepository
