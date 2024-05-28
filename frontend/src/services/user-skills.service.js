import axios from 'axios';

const API_URL = 'http://localhost:8081/user-skills';

class UserSkillsService {
  createUserSkill(userSkill) {
    return axios.post(
      API_URL + `/create`,
      {
        user_id: userSkill.userId,
        skill_id: userSkill.skillId
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

  getEntrepreneurSkills(id, page) {
    return axios.get(
      API_URL + `?entrepreneur-id=${id}&page=${page}`
    )
  }
}

export default new UserSkillsService();
