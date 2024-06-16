class UserModel {
  constructor(id, username, fullName, city, gender, birthday) {
    this.fullName = fullName;
    this.city = city;
    this.gender = gender;
    this.birthday = birthday;
    this.id = id;
    this.username = username;
  }
}

export default UserModel;