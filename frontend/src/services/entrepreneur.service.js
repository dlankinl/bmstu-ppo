import axios from 'axios';
import authHeader from './auth-header';

const API_URL = 'http://localhost:8081/entrepreneurs';

class EntrepreneurService {
  getEntrepreneurs(page) {
    // return axios.get(API_URL + `?page=${page}`);
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
    // return axios.get(API_URL + `/${id}/rating`)
    return axios.get(API_URL + `/${id}/rating`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user')}`
        }
    })
  }

  getEntrepreneurDetails(id) {
    // return axios.get(API_URL + `/${id}`)
    return axios.get(API_URL + `/${id}`, {
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('user')}`
        }
    })
  }

  // updateEntrepreneur(id, user) {
  //   return axios.patch(API_URL + `/${id}/update`, {
  //     headers: {
  //       Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
  //     },
  //     data: {
  //       id: user.id,
  //       username: user.username,
  //       full_name: user.fullName,
  //       birthday: user.birthday,
  //       gender: user.gender,
  //       city: user.city
  //     }
  //   })
  // }

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

//   getUserBoard() {
//     return axios.get(API_URL + 'user', { headers: authHeader() });
//   }

//   getModeratorBoard() {
//     return axios.get(API_URL + 'mod', { headers: authHeader() });
//   }

//   getAdminBoard() {
//     return axios.get(API_URL + 'admin', { headers: authHeader() });
//   }
}

export default new EntrepreneurService();
