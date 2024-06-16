class CompanyModel {
  constructor(id, name, activity_field, city, owner_id) { 
    this.id = id;
    this.name = name;
    this.activity_field = activity_field;
    this.city = city;
    this.owner_id = owner_id;
  }
}

export default CompanyModel;