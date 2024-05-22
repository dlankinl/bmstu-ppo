<template>
    <div>
      <template v-if="entrepreneur.full_name">
        <h1>Информация о предпринимателе</h1>
        <div>
          <h2>Полное имя: {{ entrepreneur.full_name }}</h2>
          <p>Дата рождения: {{ formatBirthday(entrepreneur.birthday) }}</p>
          <p>Город: {{ entrepreneur.city }}</p>
          <p>Пол: {{ formatGender(entrepreneur.gender) }}</p>
          <p>Рейтинг: {{ entrepreneur.rating }}</p>
        </div>
      </template>
      <template v-else>
          <h1>Информация о данном предпринимателе не заполнена.</h1>
      </template>
    </div>
  </template>
  
  <script>
  import EntrepreneurService from '../services/entrepreneur.service'
  
  export default {
    name: 'EntrepreneurPage',
    data() {
      return {
        entrepreneur: {}
      }
    },
    created() {
      this.fetchEntrepreneurDetails()
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
      }
    }
  }
  </script>