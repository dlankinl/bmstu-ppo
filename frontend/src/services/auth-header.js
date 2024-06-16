// export default function isAuth() {
//   const loggedIn = localStorage.getItem('user');

//   if (loggedIn) {
//     return true;
//   } else {
//     return false;
//   }
// }

import VueJwtDecode from 'vue-jwt-decode';

class Utils {
  isAuth() {
    const loggedIn = localStorage.getItem('user');

    if (loggedIn) {
      return true;
    } else {
      return false;
    }
  }

  decodeJWT() {
    const token = localStorage.getItem('user').replace(/"/g, '');
    return VueJwtDecode.decode(token);
  }

  getUserIdJWT() {
    const token = localStorage.getItem('user').replace(/"/g, '')
    return VueJwtDecode.decode(token).sub;
  }

  getUserRoleJWT() {
    const token = localStorage.getItem('user').replace(/"/g, '')
    return VueJwtDecode.decode(token).role;
  }
}

export default new Utils();