<template>
    <div>
      <template v-if="entrepreneur.full_name">
        <h1>Информация о предпринимателе</h1>
        <div>
          <h2>Полное имя: {{ entrepreneur.full_name }}</h2>
          <p>Дата рождения: {{ formatBirthday(entrepreneur.birthday) }}</p>
          <p>Город: {{ entrepreneur.city }}</p>
          <p>Пол: {{ formatGender(entrepreneur.gender) }}</p>
          <p>Рейтинг: {{ (5 * entrepreneur.rating).toFixed(1) }}</p>
        </div>
      </template>
      <template v-else>
          <h1>Информация о данном предпринимателе не заполнена.</h1>
      </template>
    </div>
    <div class="contacts">
        <Accordion :multiple="true">
          <AccordionTab header="Контактные данные">
            <p class="m-0">
              <div v-if="role=='guest'">Войдите в аккаунт, чтобы увидеть список средств связи.</div>
              <DataTable v-else :value="contacts" tableStyle="min-width: 30rem">
                <Column field="name" header="Название"></Column>
                <Column field="value" header="Значение"></Column>
              </DataTable>
            </p>
          </AccordionTab>
          <AccordionTab header="Финансовые показатели">
              <p class="m-0">
                <div v-if="role=='guest'">
                  Войдите в аккаунт, чтобы увидеть финансовые показатели предпринимателя.
                </div>
                <div v-else>
                  <p>Выручка: {{ financials.revenue }}</p>
                  <p>Расходы: {{ financials.costs }}</p>
                  <p>Прибыль: {{ financials.profit }}</p>
                  <p>Налоги: {{ financials.taxes }}</p>
                  <p>Налоговая нагрузка: {{ (financials.taxLoad * 100).toFixed(2) }}%</p>
                </div>               
              </p>
          </AccordionTab>
        </Accordion>
        <ButtonGroup>
          <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/companies`">
            <Button 
              icon="pi pi-building"
              label="Компании"
              class="flex-auto md:flex-initial white-space-nowrap"
            >
            </Button>
          </RouterLink>
          <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/skills`">
            <Button 
              icon="pi pi-building"
              label="Навыки"
              class="flex-auto md:flex-initial white-space-nowrap"
            >
            </Button>
          </RouterLink>
          <template v-if="role=='admin'">
            <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/update`">
              <Button 
                icon="pi pi-cog"
                label="Обновить информацию"
                class="flex-auto md:flex-initial white-space-nowrap"
              >
              </Button>
            </RouterLink>
          </template>
        </ButtonGroup>
      </div>
  </template>
  
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