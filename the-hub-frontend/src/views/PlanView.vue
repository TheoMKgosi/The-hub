<script setup lang="ts">
import { reactive, ref } from "vue";
import { RouterLink } from "vue-router";
import { useFetch } from '@vueuse/core'
import TheGoals from '../components/TheGoals.vue'
import TheTasks from '../components/TheTasks.vue'
import TheCalendar from '../components/TheCalendar.vue'

interface FormData {
  title: string
  description: string
}

const formData = reactive({
  title: '',
  description: ''
})

const currentSection = ref('tasks')

const url = "http://localhost:8080/goals"
const submitForm = async () => {
  const { data, pending, error, refresh } = await useFetch(url).post(formData).json()
  goal.value = data.value
}
</script>

<template>
  <div class="p-5">
    <h1 class="font-bold text-center text-6xl">Planning</h1>
  </div>

  <nav class="p-4 text-center text-2xl bg-gray-400">
    <button class="p-4 hover:text-blue-900" @click="currentSection = 'tasks'">Task</button>
    <button class="p-4 hover:text-blue-900" @click="currentSection = 'goals'">Goal</button>
    <button class="p-4 hover:text-blue-900" @click="currentSection = 'calendar'">Calendar</button>
  </nav>
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

  <section v-show="currentSection === 'tasks'">
    <TheTasks ></TheTasks>
  </section>

  <section v-show="currentSection === 'goals'">
    <p>Goals</p>
  </section>

  <section v-show="currentSection === 'calendar'">
    <TheCalendar >></TheCalendar>
  </section>


</template>
