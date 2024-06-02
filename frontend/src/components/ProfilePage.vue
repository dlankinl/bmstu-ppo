<template>
  <div>
    <template v-if="entrepreneur.full_name">
      <h1>Ваш профиль</h1>
      <div>
        <h2>Полное имя: {{ entrepreneur.full_name }}</h2>
        <p>Дата рождения: {{ formatBirthday(entrepreneur.birthday) }}</p>
        <p>Город: {{ entrepreneur.city }}</p>
        <p>Пол: {{ formatGender(entrepreneur.gender) }}</p>
        <p>Рейтинг: {{ (5 * entrepreneur.rating).toFixed(1) }}</p>
      </div>
    </template>
    <template v-else>
        <h1>Ваш профиль не заполнен. Обратитесь к администратору.</h1>
    </template>
  </div>
  <div v-if="entrepreneur.full_name" class="settings">
    <ButtonGroup>
      <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/contacts`"><Button label="Средства связи" icon="pi pi-address-book"></Button></RouterLink>
      <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/skills`"><Button label="Навыки" icon="pi pi-bolt"></Button></RouterLink>
      <RouterLink :to="`/entrepreneurs/${entrepreneur.id}/companies`"><Button label="Компании" icon="pi pi-building"></Button></RouterLink>
    </ButtonGroup>
  </div>
  <div v-if="role=='admin'">
    <ButtonGroup>
      <RouterLink :to="`/activity-fields`"><Button label="Все сферы деятельности" icon="pi pi-address-book"></Button></RouterLink>
      <RouterLink :to="`/skills`"><Button label="Все навыки" icon="pi pi-bolt"></Button></RouterLink>
    </ButtonGroup>
  </div>
</template>

<script>
import EntrepreneurService from '../services/entrepreneur.service'
import ButtonGroup from 'primevue/buttongroup';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Utils from '../services/auth-header';

export default {
  name: 'ProfilePage',
  components: {
    Button,
    Column,
    ButtonGroup
  },
  data() {
    return {
      entrepreneur: {},
      isAuthValue: null,
      role: 'guest'
    }
  },
  created() {
    if (Utils.isAuth()) {
      this.role = Utils.getUserRoleJWT();
    }
    this.fetchEntrepreneurDetails()
  },
  methods: {
    fetchEntrepreneurDetails() {
      const entrepreneurId = Utils.getUserIdJWT();
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
    isAuth() {
      this.isAuthValue = Utils.isAuth();
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
  }
}
</script>