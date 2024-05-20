import axios from 'axios';

const API_URL = 'http://localhost:8081/';

class AuthService {
  login(user) {
    return axios
      .post(API_URL + 'login', {
        login: user.username,
        password: user.password
      })
      .then(response => {
        console.log(response.data.data.token)
        if (response.data.data.token) {
          localStorage.setItem('user', JSON.stringify(response.data.data.token));
        }

        return response.data.data;
      });
  }

  logout() {
    localStorage.removeItem('user');
  }

  register(user) {
    return axios.post(API_URL + 'signup', {
      login: user.username,
      password: user.password
    });
  }
}

export default new AuthService();
