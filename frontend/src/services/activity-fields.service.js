import axios from 'axios';

const API_URL = 'http://localhost:8081/activity_fields';

class ActivityFieldsService {
  getActivityField(id) {
    return axios.get(API_URL + `/${id}`)
  }

  getFields() {
    return axios.get(API_URL)
  }
}

export default new ActivityFieldsService();
