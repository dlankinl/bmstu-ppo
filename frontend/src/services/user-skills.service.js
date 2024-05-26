import axios from 'axios';

const API_URL = 'http://localhost:8081/user-skills';

class UserSkillsService {
  createUserSkill(userSkill) {
    return axios.post(
      API_URL + `/create`,
      {
        user_id: userSkill.user_id,
        skill_id: userSkill.skill_id
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
}

export default new UserSkillsService();
