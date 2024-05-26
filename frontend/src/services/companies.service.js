import axios from 'axios';

const API_URL = 'http://localhost:8081/companies';

class CompaniesService {
  getEntrepreneursCompanies(id, page) {
    return axios.get(API_URL + `?entrepreneur-id=${id}&page=${page}`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
    })
  }
}

export default new CompaniesService();
