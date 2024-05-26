import axios from 'axios';

const API_URL = 'http://localhost:8081/contacts';

class ContactsService {
  getByOwnerId(id) {
    return axios.get(API_URL + `?entrepreneur-id=${id}`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
    })
  }
}

export default new ContactsService();
