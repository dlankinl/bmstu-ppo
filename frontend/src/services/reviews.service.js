import axios from 'axios';

const API_URL = 'http://localhost:8081/reviews';

class ReviewsService {
  getReview(id) {
    return axios.get(
      API_URL + `/${id}`
    )
  }

  getAuthorReviews(page) {
    return axios.get(
      API_URL + `/my?page=${page}`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  getEntrepreneurReviews(id, page) {
    return axios.get(
      API_URL + `?entrepreneur-id=${id}&page=${page}`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  createReview(rev) {
    return axios.post(
      API_URL + `/create`,
      {
        pros: rev.pros,
        cons: rev.cons,
        rating: rev.rating,
        description: rev.description,
        target_id: rev.target_id,
        reviewer_id: rev.reviewer_id
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  deleteReview(id) {
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

export default new ReviewsService();