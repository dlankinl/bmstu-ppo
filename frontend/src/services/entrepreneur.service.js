import axios from 'axios';

const API_URL = 'http://localhost:8081/entrepreneurs';

class EntrepreneurService {
  getEntrepreneurs(page) {
    return axios.get(API_URL + `?page=${page}`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user')}`
        }
    })
  }

  getEmptyEntrepreneurs(page) {
    return axios.get(API_URL + `/empty?page=${page}`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user')}`
        }
    })
  }

  getEntrepreneurRating(id) {
    return axios.get(API_URL + `/${id}/rating`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user')}`
        }
    })
  }

  getEntrepreneurDetails(id) {
    return axios.get(API_URL + `/${id}`, {
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('user')}`
        }
    })
  }

  updateEntrepreneur(id, user) {
    return axios.patch(
      API_URL + `/${id}/update`,
      {
        id: user.id,
        username: user.username,
        full_name: user.fullName,
        birthday: user.birthday,
        gender: user.gender,
        city: user.city
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }
}

export default new EntrepreneurService();
