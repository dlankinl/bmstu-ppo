import axios from 'axios';

const API_URL = 'http://localhost:8081/skills';

class SkillsService {
  getSkill(id) {
    return axios.get(
      API_URL + `/${id}`
    )
  }

  createSkill(skill) {
    return axios.post(
      API_URL + `/create`,
      {
        name: skill.name,
        description: skill.description
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  deleteSkill(id) {
    return axios.delete(
      API_URL + `/${id}/delete`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  updateSkill(id, contact) {
    return axios.patch(
      API_URL + `/${id}/update`,
      {
        name: contact.name,
        description: skill.description
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }
}

export default new SkillsService();
