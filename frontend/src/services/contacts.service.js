import axios from 'axios';

const API_URL = 'http://localhost:8081/contacts';

class ContactsService {
  getById(id) {
    return axios.get(API_URL + `/${id}`, {
      headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
      }
  })
  }

  getByOwnerId(id) {
    return axios.get(API_URL + `?entrepreneur-id=${id}`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
    })
  }

  createContact(contact) {
    return axios.post(
      API_URL + `/create`,
      {
        name: contact.name,
        value: contact.value
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  deleteContact(id) {
    return axios.delete(
      API_URL + `/${id}/delete`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  updateContact(id, contact) {
    return axios.patch(
      API_URL + `/${id}/update`,
      {
        name: contact.name,
        value: contact.value
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }
}

export default new ContactsService();
