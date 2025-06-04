<script setup lang="ts">
import { reactive } from "vue";
import { useFetch } from '@vueuse/core'
import TheGoals from '../components/TheGoals.vue'
import TheTasks from '../components/TheTasks.vue'

interface FormData {
  title: string
  description: string
}

const formData = reactive({
  title: '',
  description: ''
})

const url = "http://localhost:8080/goals"
const submitForm = async () => {
  const {data,pending,error,refresh} = await useFetch(url).post(formData).json()
  goal.value = data.value
}
</script>

<template>
  <div>
    <h1 class="font-bold text-center text-3xl">Goals and tasks</h1>
  </div>

  <!--
  <form @submit.prevent="submitForm" class="m-4">
    <fieldset class="flex flex-col border p-4 max-w-md">
      <legend class="text-center">Goal entry</legend>
      <label>Title</label>
      <input type="text" v-model="formData.title" class="border rounded-xs max-w-sm">
      <label>Description</label>
      <input type="text" v-model="formData.description" class="border rounded-xs max-w-sm">
      <button type="submit">Create goal</button>
    </fieldset>
  </form>

  <TheGoals />
-->

  <TheTasks />

</template>

