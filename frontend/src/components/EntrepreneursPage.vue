<template>
    <div>
        <h1>Entrepreneurs</h1>
        <div v-for="entrepreneur in entrepreneurs" :key="entrepreneur.id">
        <!-- <h2>Полное имя: {{ entrepreneur.full_name }}</h2> -->
        <h2>Полное имя: <a :href="`/entrepreneurs/${entrepreneur.id}`">{{ entrepreneur.full_name }}</a></h2>
        <p>Дата рождения: {{ formatBirthday(entrepreneur.birthday) }}</p>
        <p>Город: {{ entrepreneur.city }}</p>
        <p>Пол: {{ formatGender(entrepreneur.gender)  }}</p>
        <p>Рейтинг: {{ entrepreneur.rating }}</p>
        </div>
        <div>
        <button @click="prevPage" :disabled="currentPage === 1">Previous</button>
        <button @click="nextPage">Next</button>
        </div>
    </div>
</template>
  
  <script>
  import EntrepreneurService from '../services/entrepreneur.service';
  
  export default {
    name: 'EntrepreneursPage',
    data() {
      return {
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
      }
    }
  }
  </script>