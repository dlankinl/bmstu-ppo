class ReviewModel {
  constructor(id, reviewer_id, target_id, pros, cons, description, rating) {
    this.id = id;
    this.reviewer_id = reviewer_id;
    this.target_id = target_id;
    this.pros = pros;
    this.cons = cons;
    this.description = description;
    this.rating = rating;
  }
}

export default ReviewModel;