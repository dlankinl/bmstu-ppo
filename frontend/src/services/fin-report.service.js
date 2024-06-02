import axios from 'axios';

const API_URL = 'http://localhost:8081';

class FinReportService {
  getFinReport(id) {
    return axios.get(API_URL + `/financials/${id}`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user')}`
        }
    })
  }

  createFinReport(report) {
    return axios.post(
      API_URL + `/companies/${report.company_id}/financials/create`,
      {
        company_id: report.company_id,
        revenue: report.revenue,
        costs: report.costs,
        year: report.year,
        quarter: report.quarter
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  updateFinReport(id, report) {
    return axios.patch(
      API_URL + `/financials/${id}/update`,
      {
        id: report.id,
        company_id: report.company_id,
        revenue: report.revenue,
        costs: report.costs,
        year: report.year,
        quarter: report.quarter
      },
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
      }
    )
  }

  getCompanyReport(id, startYear, endYear, startQuarter, endQuarter) {
    return axios.get(API_URL + `/companies/${id}/financials/${startYear}_${startQuarter}-${endYear}_${endQuarter}`, {
      headers: {
          Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
      }
    })
  }

  getLastYearReport(id) {
    return axios.get(API_URL + `/financials?entrepreneur-id=${id}`, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('user').replace(/"/g, '')}`
        }
    })
  }
}

export default new FinReportService();
