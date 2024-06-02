<template>
  <div class="entrepreneur-details">
    <template v-if="entrepreneur.full_name">
      <div class="header">
        <h1>Информация о предпринимателе</h1>
        <div class="rating">
          <i class="pi pi-star"></i>
          {{ (5 * entrepreneur.rating).toFixed(1) }}
        </div>
      </div>
      <div class="info-card">
        <h2>{{ entrepreneur.full_name }}</h2>
        <p><i class="pi pi-calendar"></i> {{ formatBirthday(entrepreneur.birthday) }}</p>
        <p><i class="pi pi-map-marker"></i> {{ entrepreneur.city }}</p>
        <p><i class="pi pi-venus-mars"></i> {{ formatGender(entrepreneur.gender) }}</p>
      </div>
    </template>
    <template v-else>
      <h1 class="no-info">Информация о данном предпринимателе не заполнена.</h1>
    </template>
    <div class="accordion-container">
      <Accordion :multiple="true">
        <AccordionTab header="Контактные данные">
          <p class="m-0">
            <div v-if="role == 'guest'" class="guest-message">
              Войдите в аккаунт, чтобы увидеть список средств связи.
            </div>
            <DataTable v-else :value="contacts" tableStyle="min-width: 30rem">
              <Column field="name" header="Название"></Column>
              <Column field="value" header="Значение"></Column>
            </DataTable>
          </p>
        </AccordionTab>
        <AccordionTab header="Финансовые показатели">
          <p class="m-0">
            <div v-if="role == 'guest'" class="guest-message">
              Войдите в аккаунт, чтобы увидеть финансовые показатели предпринимателя.
            </div>
            <div v-else class="financials">
              <p><span class="label">Выручка:</span> {{ financials.revenue }}</p>
              <p><span class="label">Расходы:</span> {{ financials.costs }}</p>
              <p><span class="label">Прибыль:</span> {{ financials.profit }}</p>
              <p><span class="label">Налоги:</span> {{ financials.taxes }}</p>
              <p><span class="label">Налоговая нагрузка:</span> {{ (financials.taxLoad * 100).toFixed(2) }}%</p>
            </div>
          </p>
        </AccordionTab>
      </Accordion>
    </div>
    <div class="button-group">
      <ButtonGroup>
        <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/companies`">
          <Button icon="pi pi-building" label="Компании" class="button"></Button>
        </RouterLink>
        <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/skills`">
          <Button icon="pi pi-building" label="Навыки" class="button"></Button>
        </RouterLink>
        <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/reviews`">
          <Button icon="pi pi-bookmark" label="Отзывы" class="button"></Button>
        </RouterLink>
        <template v-if="role == 'admin'">
          <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/update`">
            <Button icon="pi pi-cog" label="Обновить информацию" class="button"></Button>
          </RouterLink>
        </template>
      </ButtonGroup>
    </div>
  </div>
</template>

<style scoped>
.entrepreneur-details {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f8f8f8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.rating {
  background-color: #fff;
  padding: 5px 10px;
  border-radius: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 5px;
}

.info-card {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.info-card h2 {
  margin-top: 0;
}

.info-card p {
  margin-bottom: 10px;
}

.info-card p i {
  margin-right: 5px;
}

.no-info {
  text-align: center;
  color: #888;
}

.accordion-container {
  margin-bottom: 20px;
}

.guest-message {
  text-align: center;
  color: #888;
  padding: 20px;
}

.financials {
  padding: 20px;
}

.financials p {
  margin-bottom: 10px;
}

.financials .label {
  font-weight: bold;
}

.button-group {
  display: flex;
  justify-content: center;
  gap: 10px;
}

.button {
  flex: 1;
}

@media (max-width: 768px) {
  .button-group {
    flex-wrap: wrap;
  }

  .button {
    flex: 0 0 calc(50% - 5px);
  }
}
</style>

  
<script>
  import EntrepreneurService from '../services/entrepreneur.service'
  import ContactsService from '../services/contacts.service';
  import FinReportService from '../services/fin-report.service';
  import DataTable from 'primevue/datatable';
  import Column from 'primevue/column';
  import Button from 'primevue/button';
  import Accordion from 'primevue/accordion';
  import AccordionTab from 'primevue/accordiontab';
  import Utils from '../services/auth-header';
  
  export default {
    name: 'EntrepreneurPage',
    components: {
      Button,
      Accordion,
      AccordionTab,
      DataTable,
      Column
    },
    data() {
      return {
        entrepreneur: {},
        contacts: {},
        financials: {"revenue": 0.0, "profit": 0.0, "costs": 0.0, "taxes": 0.0, "taxLoad": 0.0},
        companies: {},
        currentPage: 1,
        rows: 3,
        totalPages: 0,
        isAuthValue: null,
        role: 'guest'
      }
    },
    created() {
      this.isAuth();
      if (this.isAuthValue) {
        this.role = Utils.getUserRoleJWT();
      }
      this.fetchEntrepreneurDetails()
      this.fetchContacts()
      this.fetchFinancials()
    },
    methods: {
      fetchEntrepreneurDetails() {
        const entrepreneurId = this.$route.params.id
        EntrepreneurService.getEntrepreneurDetails(entrepreneurId)
          .then(response => {
            const entrepreneur = response.data.data.entrepreneur
            if (entrepreneur.full_name) {
                this.fetchRating(entrepreneur)
            } else {
                this.entrepreneur = entrepreneur
            }
          })
          .catch(error => {
            console.error('Error fetching entrepreneur details:', error)
          })
      },
      formatBirthday(birthday) {
        const date = new Date(birthday)
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
      },
      formatGender(gender) {
        return gender === 'm' ? 'мужской' : 'женский'
      },
      fetchRating(entrepreneur) {
        EntrepreneurService.getEntrepreneurRating(entrepreneur.id)
          .then(response => {
            entrepreneur.rating = response.data.data.rating
            this.entrepreneur = entrepreneur
          })
          .catch(error => {
            console.error('Error fetching entrepreneur rating:', error)
          })
      },
      fetchContacts() {
        const id = this.$route.params.id;
        if (this.isAuthValue) {
          ContactsService.getByOwnerId(id)
          .then(response => {
            this.contacts = response.data.data.contacts;
          })
          .catch(error => {
            console.error('Ошибка получения средств связи:', error)
          })
        }
      },
      fetchFinancials() {
        const id = this.$route.params.id;
        if (this.isAuthValue) {
          FinReportService.getLastYearReport(id)
            .then(response => {
              this.financials = response.data.data;
            })
            .catch(error => {
              console.error('Ошибка получения прошлогоднего отчета:', error)
            })
        }
      },
      isAuth() {
        this.isAuthValue = Utils.isAuth();
      },
    }
  }
  </script>