<script setup lang="ts">
const taskStore = useTaskStore()
</script>

<template>
  <main>
    <TheWelcome />
    <!-- Main Dashboard -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Overview Stats -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
              <span class="text-green-600 font-bold">✓</span>
            </div>
            <div class="ml-4 overflow-auto">
              <p class="text-sm font-medium text-gray-600">Tasks completed</p>
              <p class="text-2xl font-semibold text-gray-900">{{ taskStore.completedTasks.length }}/ {{
                taskStore.tasks.length }}
              </p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
              <span class="text-blue-600 font-bold">⏰</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Hours Today</p>
              <p class="text-2xl font-semibold text-gray-900">6.5h</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-purple-100 rounded-lg flex items-center justify-center">
              <span class="text-purple-600 font-bold">📚</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Courses Active</p>
              <p class="text-2xl font-semibold text-gray-900">3</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-emerald-100 rounded-lg flex items-center justify-center">
              <span class="text-emerald-600 font-bold">$</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Budget Left</p>
              <p class="text-2xl font-semibold text-gray-900">$2,340</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Management Sections Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div class="bg-white rounded-lg shadow">
          <DashboardTask />
        </div>
        <!-- Time Management -->
        <div class="bg-white rounded-lg shadow">
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 flex items-center">
                <span class="text-blue-600 mr-2">⏰</span>
                Time Management
              </h2>
              <button class="text-gray-400 hover:text-gray-600">
                <span class="text-lg">📅</span>
              </button>
            </div>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="flex items-center justify-between p-3 bg-blue-50 rounded-lg border-l-4 border-blue-500">
                <div>
                  <p class="text-sm font-medium text-gray-900">Team Standup</p>
                  <p class="text-xs text-gray-600">9:00 AM - 9:30 AM</p>
                </div>
                <span class="text-xs text-blue-600 font-medium">Now</span>
              </div>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p class="text-sm font-medium text-gray-900">Design Review</p>
                  <p class="text-xs text-gray-600">11:00 AM - 12:00 PM</p>
                </div>
                <span class="text-xs text-gray-500">Upcoming</span>
              </div>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p class="text-sm font-medium text-gray-900">Client Presentation</p>
                  <p class="text-xs text-gray-600">2:00 PM - 3:00 PM</p>
                </div>
                <span class="text-xs text-gray-500">Upcoming</span>
              </div>
            </div>
            <div class="mt-6">
              <h3 class="text-sm font-medium text-gray-900 mb-3">Time Tracking</h3>
              <div class="space-y-2">
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Development</span>
                  <span class="text-sm font-medium text-gray-900">4.5h</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Meetings</span>
                  <span class="text-sm font-medium text-gray-900">1.5h</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Planning</span>
                  <span class="text-sm font-medium text-gray-900">0.5h</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Learning Management -->
        <div class="bg-white rounded-lg shadow">
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 flex items-center">
                <span class="text-purple-600 mr-2">📚</span>
                Learning Management
              </h2>
              <button class="text-gray-400 hover:text-gray-600">
                <span class="text-lg">🎯</span>
              </button>
            </div>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="p-4 bg-purple-50 rounded-lg">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-gray-900">Advanced React Patterns</h3>
                  <span class="text-xs text-purple-600 font-medium">85%</span>
                </div>
                <div class="w-full bg-purple-200 rounded-full h-2 mb-2">
                  <div class="bg-purple-600 h-2 rounded-full" style="width: 85%"></div>
                </div>
                <p class="text-xs text-gray-600">12/14 lessons completed</p>
              </div>
              <div class="p-4 bg-blue-50 rounded-lg">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-gray-900">TypeScript Fundamentals</h3>
                  <span class="text-xs text-blue-600 font-medium">60%</span>
                </div>
                <div class="w-full bg-blue-200 rounded-full h-2 mb-2">
                  <div class="bg-blue-600 h-2 rounded-full" style="width: 60%"></div>
                </div>
                <p class="text-xs text-gray-600">9/15 lessons completed</p>
              </div>
              <div class="p-4 bg-green-50 rounded-lg">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-gray-900">System Design Basics</h3>
                  <span class="text-xs text-green-600 font-medium">30%</span>
                </div>
                <div class="w-full bg-green-200 rounded-full h-2 mb-2">
                  <div class="bg-green-600 h-2 rounded-full" style="width: 30%"></div>
                </div>
                <p class="text-xs text-gray-600">3/10 lessons completed</p>
              </div>
            </div>
            <div class="mt-6">
              <h3 class="text-sm font-medium text-gray-900 mb-3">Study Goals</h3>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <span class="text-sm text-gray-900">Daily study time</span>
                <span class="text-sm font-medium text-green-600">2h / 2h ✓</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Financial Management -->
        <div class="bg-white rounded-lg shadow">
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 flex items-center">
                <span class="text-emerald-600 mr-2">💰</span>
                Financial Management
              </h2>
              <button class="text-gray-400 hover:text-gray-600">
                <span class="text-lg">📈</span>
              </button>
            </div>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-2 gap-4 mb-6">
              <div class="text-center p-4 bg-emerald-50 rounded-lg">
                <p class="text-sm text-gray-600">Income</p>
                <p class="text-xl font-semibold text-emerald-600">$5,240</p>
              </div>
              <div class="text-center p-4 bg-red-50 rounded-lg">
                <p class="text-sm text-gray-600">Expenses</p>
                <p class="text-xl font-semibold text-red-600">$2,900</p>
              </div>
            </div>
            <div class="space-y-3">
              <h3 class="text-sm font-medium text-gray-900">Recent Transactions</h3>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div class="flex items-center">
                  <div class="w-8 h-8 bg-emerald-100 rounded-full flex items-center justify-center">
                    <span class="text-emerald-600 text-sm">+</span>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">Freelance Project</p>
                    <p class="text-xs text-gray-600">Dec 15, 2024</p>
                  </div>
                </div>
                <span class="text-sm font-medium text-emerald-600">+$1,200</span>
              </div>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div class="flex items-center">
                  <div class="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                    <span class="text-red-600 text-sm">-</span>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">Office Supplies</p>
                    <p class="text-xs text-gray-600">Dec 14, 2024</p>
                  </div>
                </div>
                <span class="text-sm font-medium text-red-600">-$85</span>
              </div>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div class="flex items-center">
                  <div class="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                    <span class="text-red-600 text-sm">-</span>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">Software License</p>
                    <p class="text-xs text-gray-600">Dec 12, 2024</p>
                  </div>
                </div>
                <span class="text-sm font-medium text-red-600">-$29</span>
              </div>
            </div>
            <div class="mt-6">
              <h3 class="text-sm font-medium text-gray-900 mb-3">Budget Overview</h3>
              <div class="space-y-2">
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Food & Dining</span>
                  <span class="text-sm text-gray-900">$450 / $600</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-2">
                  <div class="bg-yellow-500 h-2 rounded-full" style="width: 75%"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </main>
</template>
