<template>
  <div class="pt-4">
    <div class="welcome-card">
      <h2 class="welcome-title">Добро пожаловать в наше сообщество предпринимателей!</h2>
      <p class="welcome-text">
        Это приложение создано для людей, кто давно жаждет открыть свое дело или новое направление, но боится, что не справится один!
      </p>
      <p class="welcome-text">
        Здесь вы можете найти единомышленников, связаться с ними и превратить свою мечту в реальность. Идеальное место для нетворкинга и установления полезных связей.
      </p>
      <p class="welcome-text">
        Помните, что конфиденциальность личных данных соблюдается внимательно, а каждый новый зарегистрированный пользователь становится полноценным Предпринимателем на данной платформе только после документальной проверки у администратора. Вы можете не бояться нарваться на мошенника или человека, представляющегося другим. У нас всё надежно!
      </p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomePage'
}
</script>

<style scoped>
.welcome-card {
  background-color: #f8f8f8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.welcome-title {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin-bottom: 16px;
}

.welcome-text {
  font-size: 16px;
  color: #555;
  line-height: 1.5;
  margin-bottom: 12px;
}
</style>