<template>
  <div>
    <h1>Entrepreneurs</h1>
    <div v-for="entrepreneur in entrepreneurs" :key="entrepreneur.id">
    <template v-if="entrepreneur.full_name">
      <h2>Полное имя: <a :href="`/entrepreneurs/${entrepreneur.id}`">{{ entrepreneur.full_name }}</a></h2>
      <p>Дата рождения: {{ formatBirthday(entrepreneur.birthday) }}</p>
      <p>Город: {{ entrepreneur.city }}</p>
      <p>Пол: {{ formatGender(entrepreneur.gender)  }}</p>
      <p>Рейтинг: {{ entrepreneur.rating }}</p>
    </template>
    <template v-else>
      <h2>Имя пользователя: <a :href="`/entrepreneurs/${entrepreneur.id}`">{{ entrepreneur.username }}</a></h2>
    </template>
    </div>
    <div>
    <Button @click="prevPage" :disabled="currentPage === 1" plain text raised>Назад</Button>
    <Button @click="nextPage" plain text raised>Далее</Button>
    <!-- <b-button variant="outline-primary" @click="nextPage">Далее</b-button>
    <b-button variant="outline-primary" @click="prevPage" :disabled="currentPage === 1">Назад</b-button> -->
      <!-- <Paginator 
        v-model:first="first"
        :rows="rows"
        :totalRecords="totalRecords"
        :rowsPerPageOptions="[3]"
        @page="onPageChange"
      ></Paginator> -->
    </div>
  </div>
</template>
  
  <script>
  import EntrepreneurService from '../services/entrepreneur.service';
  import Button from 'primevue/button';
  import Paginator from 'primevue/paginator';
  
  export default {
    name: 'EntrepreneursPage',
    components: {
      Button,
      Paginator
    },
    data() {
      return {
        first: 0,
        rows: 3,
        totalRecords: 0,
        entrepreneurs: [],
        currentPage: 1
      }
    },
    created() {
      this.fetchEntrepreneurs()
    },
    methods: {
      fetchEntrepreneurs() {
        EntrepreneurService.getEntrepreneurs(this.currentPage)
          .then(response => {
            const entrepreneurs = response.data.data.users
            this.numPages = response.data.data.num_pages
            this.fetchRatings(entrepreneurs)
          })
          .catch(error => {
            console.error(error)
          })
      },
      prevPage() {
        this.currentPage--
        this.fetchEntrepreneurs()
      },
      nextPage() {
        this.currentPage++
        this.fetchEntrepreneurs()
      },
      formatGender(gender) {
        return gender === 'm' ? 'мужской' : 'женский'
      },
      formatBirthday(birthday) {
        const date = new Date(birthday)
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
      },
      fetchRatings(entrepreneurs) {
        Promise.all(
          entrepreneurs.map(entrepreneur =>
            EntrepreneurService.getEntrepreneurRating(entrepreneur.id).then(response => {
              entrepreneur.rating = response.data.data.rating
            })
          )
        ).then(() => {
          this.entrepreneurs = entrepreneurs
        })
      },
      onPageChange(event) {
        this.first = event.first;
        this.rows = event.rows;
        this.currentPage = Math.floor(this.first / this.rows) + 1;
        this.fetchEntrepreneurs();
      },
      // fetchEntrepreneurs() {
      //   EntrepreneurService.getEntrepreneurs(this.currentPage)
      //     .then(response => {
      //       this.entrepreneurs = response.data.data.users;
      //       this.totalRecords = response.data.data.total_users;
      //       this.numPages = Math.ceil(this.totalRecords / this.rows);
      //       this.fetchRatings(this.entrepreneurs);
      //     })
      //     .catch(error => {
      //       console.error(error);
      //     });
      // },
    }
  }
  </script>