// eslint-disable-next-line
import axios from "axios";

// A simple data API that will be used to get the data for our
// components. On a real website, a more robust data fetching
// solution would be more appropriate.
const API = {
  all(callback) {
    axios
      .get('http://localhost:8000/api/index')
      .then((resp) => {
        callback(resp)
      })
  },
  get(name) {
    const isRepo = p => p.name === name
    return this.repositories.find(isRepo)
  },
}

export default API
