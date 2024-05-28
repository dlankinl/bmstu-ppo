<template>
  <div class="card flex flex-column md:flex-row gap-3">
    <Message v-if="message" :severity="message.severity" :life="3000">{{ message.content }}</Message>
    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-user"></i>
      </InputGroupAddon>
      <InputText v-model="name" placeholder="Название" :class="{ 'p-invalid': !name }" />
    </InputGroup>

    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-map"></i>
      </InputGroupAddon>
      <InputText v-model="description" placeholder="Описание" :class="{ 'p-invalid': !description }"/>
    </InputGroup>
    <Button @click="updateSkill">Обновить</Button>
  </div>
</template>
    
<script>
import SkillsService from '../services/skills.service';
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import SkillModel from '../models/SkillModel.js'
import Message from 'primevue/message';
import { ref } from 'vue';
  
export default {
  name: 'UpdateSkillPage',
  components: {
    Button,
    InputGroup,
    InputGroupAddon,
    InputText,
    Message
  },
  setup() {
    const name = ref('');
    const description = ref('');
    const message = ref(null);

    return {
      name,
      description,
      message
    };
  },
  created() {
    this.fillInfo();
  },
  methods: {
    fillInfo() {
      SkillsService.getSkill(this.$route.params.id)
        .then(response => {
          this.skill = response.data.data.skill
          this.name = this.skill.name;
          this.description = this.skill.description;
        })
        .catch(error => {
          this.message = { severity: 'error', content: `Произошла ошибка при получении данных: ${error.response.data.error}` }
          console.error(error)
        })
    },
    updateSkill() {
      const skill = new SkillModel(this.$route.params.id, this.name, this.description);

      SkillsService.updateSkill(this.$route.params.id, skill)
      .then(response => {
        this.message = { severity: 'success', content: 'Данные успешно обновлены' }
      })
      .catch(error => {
        this.message = { severity: 'error', content: `Произошла ошибка при обновлении данных: ${error.response.data.error}` }
        console.error(error)
      })
    }
  }
}

</script>