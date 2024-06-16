class FinReportModel {
    constructor(id, company_id, revenue, costs, year, quarter) { 
      this.id = id;
      this.company_id = company_id;
      this.revenue = revenue;
      this.costs = costs;
      this.year = year;
      this.quarter = quarter;
    }
  }
  
  export default FinReportModel;