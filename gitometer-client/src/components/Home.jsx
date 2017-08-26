// eslint-disable-next-line
import React from 'react'
import RepositorySummary from './RepositorySummary'
import API from '../api'

const h1 = {
  textAlign: 'center',
  marginTop: '30px',
  fontWeight: '600',
  lineHeight: '1.1',
  color: '#334152',
  fontSize: '60px',
}

const Home = () =>
  (<div>
    <div className="section-heading">
      <img src={'img/logo.png'} className="center" alt="" style={{ display: 'block' }} />
      <h1 style={h1}>Gitometer</h1>
    </div>
    <div className="padded-lg">
      <div className="projects-list">
        {API.all().map(r =>
          <RepositorySummary key={r.name} name={r.name} stars={r.stars} id={r.id} />,
        )}
      </div>
    </div>
  </div>)

export default Home
