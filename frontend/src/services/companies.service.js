import axios from 'axios';

const API_URL = 'http://localhost:8081/companies';

class CompaniesService {
  getEntrepreneursCompanies(id, page) {
    return axios.get(API_URL + `?entrepreneur-id=${id}&page=${page}`)
  }

  getCompanyDetails(id) {
    return axios.get(API_URL + `/${id}`)
  }

  updateCompany(id, company) {
    return axios.patch(
      API_URL + `/${id}/update`,
      {
        name: company.name,
        activity_field_id: company.activity_field,
        city: company.city
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  deleteCompany(id) {
    return axios.delete(
      API_URL + `/${id}/delete`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  createCompany(company) {
    return axios.post(
      API_URL + `/create`,
      {
        name: company.name,
        activity_field_id: company.activity_field,
        city: company.city
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }
}

export default new CompaniesService();
