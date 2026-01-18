<template>
  <div>
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

  </div>
</template>
