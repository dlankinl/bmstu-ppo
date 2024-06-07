<template>
  <div class="card">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <div class="form-group">
      <InputGroup>
        <InputNumber v-model="startYear" :useGrouping="false" :min="1900" :max="`${new Date().getFullYear()}`" placeholder="Начальный год" :class="{ 'p-invalid': !startYear }"/>
      </InputGroup>
    </div>

    <div class="form-group">
      <InputGroup>
        <InputNumber v-model="startQuarter" :min="1" :max="4" placeholder="Начальный квартал" :class="{ 'p-invalid': !startQuarter }"/>
      </InputGroup>
    </div>

    <div class="form-group">
      <InputGroup>
        <InputNumber v-model="endYear" :useGrouping="false" :min="1900" :max="`${new Date().getFullYear()}`" placeholder="Конечный год" :class="{ 'p-invalid': !endYear }"/>
      </InputGroup>
    </div>

    <div class="form-group">
      <InputGroup>
        <InputNumber v-model="endQuarter" :min="1" :max="4" placeholder="Конечный квартал" :class="{ 'p-invalid': !endQuarter }"/>
      </InputGroup>
    </div>

    <Button @click="fetchReport" class="fetch-button">Получить отчёт</Button>

    <template v-if="financialsAvailable">
      <div class="financials">
        <p>Выручка: {{ revenue }}</p>
        <p>Расходы: {{ costs }}</p>
        <p>Прибыль: {{ profit }}</p>
      </div>
    </template>
  </div>
</template>

<style scoped>
.card {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f8f8f8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 20px;
}

.fetch-button {
  width: 100%;
}

.financials {
  margin-top: 20px;
}
</style>
  
<script>
import FinReportService from '../services/fin-report.service';
import Utils from '../services/auth-header';
import Button from 'primevue/button';
import DataView from 'primevue/dataview';
import Message from 'primevue/message';
import InputGroup from 'primevue/inputgroup';
import InputNumber from 'primevue/inputnumber';
import {ref} from 'vue';

export default {
  name: 'CompanyFinancialsPage',
  components: {
    Button,
    DataView,
    Message,
    InputNumber,
    InputGroup
  },
  data() {
    return {
      entId: null,
      message: null,
      financials: null,
      startYear: ref(null),
      startQuarter: ref(null),
      endYear: ref(null),
      endQuarter: ref(null),
      financialsAvailable: false,
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.visitorId = Utils.getUserIdJWT();
    }
    this.fetchReport();
  },
  methods: {
    fetchReport() {
      this.compId = this.$route.params.id;
      if (this.startYear != null && this.startQuarter != null && this.endYear != null & this.endQuarter != null) {
        FinReportService.getCompanyReport(this.compId, this.startYear, this.endYear, this.startQuarter, this.endQuarter)
        .then(response => {
          this.report = response.data.data;
          this.revenue = this.report.revenue;
          this.costs = this.report.costs;
          this.profit = this.report.profit;
          this.financialsAvailable = true;
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при получении отчета компании: ${error.response.data.error}` }
          console.error('Ошибка получения отчета компании: ', error)
        })
      }
    } 
  }
}
</script>