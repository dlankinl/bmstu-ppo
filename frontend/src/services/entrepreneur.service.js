import axios from 'axios';
import authHeader from './auth-header';

const API_URL = 'http://localhost:8081/entrepreneurs';

class EntrepreneurService {
  getEntrepreneurs(page) {
    return axios.get(API_URL + `?page=${page}`);
  }

  getEntrepreneurRating(id) {
    return axios.get(API_URL + `/${id}/rating`)
  }

  getEntrepreneurDetails(id) {
    return axios.get(API_URL + `/${id}`)
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
