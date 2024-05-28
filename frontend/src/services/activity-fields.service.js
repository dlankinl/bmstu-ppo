import axios from 'axios';

const API_URL = 'http://localhost:8081/activity_fields';

class ActivityFieldsService {
  getActivityField(id) {
    return axios.get(API_URL + `/${id}`)
  }

  getFields() {
    return axios.get(API_URL)
  }

  getPaginatedFields(page) {
    return axios.get(API_URL + `?page=${page}`)
  }

  createField(field) {
    return axios.post(
      API_URL + `/create`,
      {
        name: field.name,
        description: field.description,
        cost: field.cost
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  updateField(id, field) {
    return axios.patch(
      API_URL + `/${id}/update`,
      {
        name: field.name,
        description: field.description,
        cost: field.cost
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }
}

export default new ActivityFieldsService();
