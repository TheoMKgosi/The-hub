<script setup lang="ts">
import { useMouse, useElementBounding, useIntersectionObserver } from '@vueuse/core'

definePageMeta({
  layout: false
})

const auth = useAuthStore()

onMounted(() => {
  if (auth.isLoggedIn) {
    navigateTo('/dashboard')
  }
})

const { x: mouseX, y: mouseY } = useMouse()
const heroRef = ref<HTMLElement | null>(null)
const aboutRef = ref<HTMLElement | null>(null)
const featuresRef = ref<HTMLElement | null>(null)

const { width, height } = useElementBounding(heroRef)

const parallaxOffset = computed(() => {
  if (!width.value || !height.value) return { x: 0, y: 0 }
  const centerX = width.value / 2
  const centerY = height.value / 2
  return {
    x: ((mouseX.value - centerX) / centerX) * 15,
    y: ((mouseY.value - centerY) / centerY) * 15
  }
})

const heroVisible = ref(false)
const aboutVisible = ref(false)
const featuresVisible = ref(false)

onMounted(() => {
  setTimeout(() => heroVisible.value = true, 100)
})

useIntersectionObserver(aboutRef, ([{ isIntersecting }]) => {
  if (isIntersecting) aboutVisible.value = true
}, { threshold: 0.2 })

useIntersectionObserver(featuresRef, ([{ isIntersecting }]) => {
  if (isIntersecting) featuresVisible.value = true
}, { threshold: 0.1 })

const scrollToSection = (id: string) => {
  const element = document.getElementById(id)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}

const navigateToApp = () => {
  if (auth.isLoggedIn) {
    navigateTo('/dashboard')
  } else {
    navigateTo('/login')
  }
}

const features = [
  {
    title: 'Task Management',
    description: 'Organize your life with smart task lists, priorities, and deadlines. Stay on top of everything with intuitive categorization and progress tracking.',
    color: 'primary'
  },
  {
    title: 'Finance Management',
    description: 'Take control of your finances with budget tracking, expense monitoring, and savings goals. Visualize your financial health with insightful charts.',
    color: 'secondary'
  },
  {
    title: 'Learning Management',
    description: 'Accelerate your learning with spaced repetition flashcards, structured study paths, and detailed progress analytics. Master any subject efficiently.',
    color: 'accent'
  }
]

const navScrolled = ref(false)

onMounted(() => {
  const handleScroll = () => {
    navScrolled.value = window.scrollY > 50
  }
  window.addEventListener('scroll', handleScroll, { passive: true })
  onUnmounted(() => window.removeEventListener('scroll', handleScroll))
})
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark overflow-x-hidden">
    <!-- Navigation -->
    <nav 
      class="fixed top-0 left-0 right-0 z-50 transition-all duration-300"
      :class="navScrolled ? 'bg-white/80 dark:bg-black/80 backdrop-blur-lg shadow-sm' : 'bg-transparent'"
    >
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <div class="flex items-center justify-between h-20">
          <!-- Logo -->
          <div class="flex items-center space-x-3">
            <img src="/logo.svg" alt="Project Life Ledger" class="h-10 w-10" />
            <span class="text-xl font-bold text-text-light dark:text-text-dark hidden sm:block">
              Project Life Ledger
            </span>
          </div>

          <!-- Nav Links -->
          <div class="hidden md:flex items-center space-x-8">
            <button 
              @click="scrollToSection('about')"
              class="text-text-light/70 dark:text-text-dark/70 hover:text-text-light dark:hover:text-text-dark transition-colors duration-200 text-sm font-medium"
            >
              About
            </button>
            <button 
              @click="scrollToSection('features')"
              class="text-text-light/70 dark:text-text-dark/70 hover:text-text-light dark:hover:text-text-dark transition-colors duration-200 text-sm font-medium"
            >
              Features
            </button>
          </div>

          <!-- CTA -->
          <BaseButton 
            @click="navigateToApp" 
            text="Get Started" 
            variant="primary" 
            size="md"
            class="transform hover:scale-105 transition-transform duration-200"
          />
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section 
      ref="heroRef"
      class="relative min-h-screen flex items-center pt-20 overflow-hidden"
    >
      <!-- Background Gradient -->
      <div class="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-secondary/5 dark:from-primary/10 dark:to-secondary/10" />
      
      <!-- Floating Orbs -->
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/20 rounded-full blur-3xl animate-pulse" />
      <div class="absolute bottom-1/4 right-1/4 w-80 h-80 bg-secondary/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;" />

      <div class="relative mx-auto max-w-7xl px-6 lg:px-8 w-full">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-8 items-center">
          <!-- Left: Text Content -->
          <div 
            class="text-center lg:text-left"
            :class="heroVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
            style="transition: all 0.8s cubic-bezier(0.16, 1, 0.3, 1)"
          >
            <!-- Badge -->
            <div 
              class="inline-flex items-center px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium mb-8"
              :class="heroVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-95'"
              style="transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1) 0.1s"
            >
              <span class="w-2 h-2 rounded-full bg-primary mr-2 animate-pulse" />
              Now available for everyone
            </div>

            <!-- Headline -->
            <h1 
              class="text-4xl sm:text-5xl lg:text-6xl font-bold text-text-light dark:text-text-dark leading-tight mb-6"
              :class="heroVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
              style="transition: all 0.7s cubic-bezier(0.16, 1, 0.3, 1) 0.2s"
            >
              Organize Your Life,
              <span class="text-primary">One Task at a Time</span>
            </h1>

            <!-- Description -->
            <p 
              class="text-lg text-text-light/70 dark:text-text-dark/70 max-w-xl mx-auto lg:mx-0 mb-8 leading-relaxed"
              :class="heroVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
              style="transition: all 0.7s cubic-bezier(0.16, 1, 0.3, 1) 0.3s"
            >
              The all-in-one productivity platform for individuals. Master your tasks, manage your finances, and accelerate your learning with intelligent tools designed for modern life.
            </p>

            <!-- CTA -->
            <div 
              class="flex flex-col sm:flex-row items-center gap-4 justify-center lg:justify-start"
              :class="heroVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
              style="transition: all 0.7s cubic-bezier(0.16, 1, 0.3, 1) 0.4s"
            >
              <BaseButton 
                @click="navigateToApp" 
                text="Get Started Free" 
                variant="primary" 
                size="lg"
                class="w-full sm:w-auto transform hover:scale-105 transition-all duration-200 shadow-lg shadow-primary/25"
              />
            </div>

            <!-- Trust Badge -->
            <p 
              class="mt-6 text-sm text-text-light/50 dark:text-text-dark/50"
              :class="heroVisible ? 'opacity-100' : 'opacity-0'"
              style="transition: opacity 0.7s ease 0.6s"
            >
              Free forever • No credit card required
            </p>
          </div>

          <!-- Right: Isometric Illustration -->
          <div 
            class="relative h-[400px] lg:h-[500px]"
            :class="heroVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-95'"
            style="transition: all 1s cubic-bezier(0.16, 1, 0.3, 1) 0.3s"
          >
            <!-- Isometric Scene Container -->
            <div class="absolute inset-0 flex items-center justify-center">
              <div 
                class="isometric-scene"
                :style="{
                  transform: `rotateX(60deg) rotateZ(-45deg) translateX(${parallaxOffset.x}px) translateY(${parallaxOffset.y}px)`
                }"
              >
                <!-- Task Card -->
                <div 
                  class="iso-element iso-task floating"
                  style="animation-delay: 0s;"
                >
                  <div class="iso-card">
                    <div class="iso-card-face iso-card-front">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-4 h-4 rounded bg-primary/20" />
                        <div class="w-20 h-2 rounded bg-text-light/20" />
                      </div>
                      <div class="space-y-2">
                        <div class="w-full h-2 rounded bg-text-light/10" />
                        <div class="w-3/4 h-2 rounded bg-text-light/10" />
                      </div>
                    </div>
                    <div class="iso-card-face iso-card-side" />
                    <div class="iso-card-face iso-card-top" />
                  </div>
                </div>

                <!-- Finance Dashboard -->
                <div 
                  class="iso-element iso-finance floating"
                  style="animation-delay: 0.5s;"
                >
                  <div class="iso-card bg-secondary/10">
                    <div class="iso-card-face iso-card-front p-4">
                      <div class="text-xs font-bold text-secondary mb-2">$12,450</div>
                      <div class="flex items-end gap-1 h-12">
                        <div class="w-3 bg-secondary/40 rounded-t" style="height: 40%;" />
                        <div class="w-3 bg-secondary/60 rounded-t" style="height: 70%;" />
                        <div class="w-3 bg-secondary/80 rounded-t" style="height: 55%;" />
                        <div class="w-3 bg-secondary rounded-t" style="height: 90%;" />
                      </div>
                    </div>
                    <div class="iso-card-face iso-card-side bg-secondary/20" />
                    <div class="iso-card-face iso-card-top bg-secondary/30" />
                  </div>
                </div>

                <!-- Learning Stack -->
                <div 
                  class="iso-element iso-learning floating"
                  style="animation-delay: 1s;"
                >
                  <div class="iso-stack">
                    <div class="iso-flashcard" style="transform: translateZ(0px);">
                      <div class="text-xs text-accent font-bold">Q: What is...</div>
                    </div>
                    <div class="iso-flashcard" style="transform: translateZ(8px);">
                      <div class="text-xs text-accent/70">A: It is...</div>
                    </div>
                    <div class="iso-flashcard" style="transform: translateZ(16px);">
                      <div class="w-full h-2 rounded bg-accent/20" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- App Preview Section -->
    <section class="relative z-10 -mt-20 pb-24">
      <div class="mx-auto max-w-6xl px-6 lg:px-8">
        <div 
          class="relative rounded-2xl overflow-hidden shadow-2xl shadow-primary/10 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark"
          :class="heroVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-12'"
          style="transition: all 1s cubic-bezier(0.16, 1, 0.3, 1) 0.6s"
        >
          <!-- Browser Chrome -->
          <div class="bg-surface-light dark:bg-surface-dark border-b border-surface-light dark:border-surface-dark px-4 py-3 flex items-center gap-2">
            <div class="flex gap-1.5">
              <div class="w-3 h-3 rounded-full bg-red-400" />
              <div class="w-3 h-3 rounded-full bg-yellow-400" />
              <div class="w-3 h-3 rounded-full bg-green-400" />
            </div>
            <div class="flex-1 mx-4">
              <div class="bg-background-light dark:bg-background-dark rounded-md px-3 py-1.5 text-xs text-text-light/50 dark:text-text-dark/50 text-center">
                project-life-ledger.app/dashboard
              </div>
            </div>
          </div>
          
          <!-- Dashboard Preview -->
          <div class="p-6 bg-background-light dark:bg-background-dark">
            <div class="grid grid-cols-12 gap-4">
              <!-- Sidebar -->
              <div class="col-span-3 bg-surface-light dark:bg-surface-dark rounded-xl p-4">
                <div class="space-y-3">
                  <div class="flex items-center gap-2 p-2 rounded-lg bg-primary/10">
                    <div class="w-8 h-8 rounded-lg bg-primary/20" />
                    <div class="w-16 h-2 rounded bg-text-light/20" />
                  </div>
                  <div v-for="i in 4" :key="i" class="flex items-center gap-2 p-2">
                    <div class="w-8 h-8 rounded-lg bg-text-light/5" />
                    <div class="w-20 h-2 rounded bg-text-light/10" />
                  </div>
                </div>
              </div>
              
              <!-- Main Content -->
              <div class="col-span-9 space-y-4">
                <!-- Stats Row -->
                <div class="grid grid-cols-3 gap-4">
                  <div v-for="i in 3" :key="i" class="bg-surface-light dark:bg-surface-dark rounded-xl p-4">
                    <div class="w-20 h-2 rounded bg-text-light/10 mb-2" />
                    <div class="w-12 h-6 rounded bg-text-light/20" />
                  </div>
                </div>
                
                <!-- Content Area -->
                <div class="bg-surface-light dark:bg-surface-dark rounded-xl p-4 h-48">
                  <div class="flex justify-between mb-4">
                    <div class="w-32 h-3 rounded bg-text-light/20" />
                    <div class="w-20 h-3 rounded bg-primary/20" />
                  </div>
                  <div class="space-y-3">
                    <div v-for="i in 4" :key="i" class="flex items-center gap-3 p-3 rounded-lg bg-background-light dark:bg-background-dark">
                      <div class="w-5 h-5 rounded border-2 border-primary/30" />
                      <div class="flex-1 h-2 rounded bg-text-light/10" />
                      <div class="w-16 h-2 rounded bg-text-light/5" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- About Section -->
    <section 
      id="about" 
      ref="aboutRef"
      class="py-24 relative overflow-hidden"
    >
      <div class="mx-auto max-w-4xl px-6 lg:px-8 text-center">
        <div
          :class="aboutVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          style="transition: all 0.8s cubic-bezier(0.16, 1, 0.3, 1)"
        >
          <!-- Section Label -->
          <span class="inline-block px-4 py-1.5 rounded-full bg-secondary/10 text-secondary text-sm font-medium mb-6">
            About
          </span>
          
          <!-- Headline -->
          <h2 class="text-3xl sm:text-4xl lg:text-5xl font-bold text-text-light dark:text-text-dark mb-6">
            Your Personal Productivity
            <span class="text-secondary">Command Center</span>
          </h2>
          
          <!-- Description -->
          <p class="text-lg text-text-light/70 dark:text-text-dark/70 max-w-2xl mx-auto leading-relaxed">
            Project Life Ledger brings together everything you need to stay organized and productive. 
            Whether you're managing daily tasks, tracking your finances, or mastering new skills, 
            our unified platform adapts to your workflow and helps you achieve your goals.
          </p>
          
          <!-- Decorative Element -->
          <div class="mt-12 flex justify-center">
            <div class="flex items-center gap-4">
              <div class="w-16 h-0.5 bg-gradient-to-r from-transparent to-primary/30" />
              <div class="w-3 h-3 rotate-45 bg-primary/20" />
              <div class="w-16 h-0.5 bg-gradient-to-l from-transparent to-secondary/30" />
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section 
      id="features" 
      ref="featuresRef"
      class="py-24 bg-surface-light/30 dark:bg-surface-dark/30"
    >
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <!-- Section Header -->
        <div 
          class="text-center mb-16"
          :class="featuresVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          style="transition: all 0.8s cubic-bezier(0.16, 1, 0.3, 1)"
        >
          <span class="inline-block px-4 py-1.5 rounded-full bg-accent/10 text-accent text-sm font-medium mb-6">
            Features
          </span>
          <h2 class="text-3xl sm:text-4xl font-bold text-text-light dark:text-text-dark mb-4">
            Everything You Need
          </h2>
          <p class="text-lg text-text-light/70 dark:text-text-dark/70 max-w-2xl mx-auto">
            Three powerful modules working together to streamline your life
          </p>
        </div>

        <!-- Features Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div
            v-for="(feature, index) in features"
            :key="feature.title"
            class="group relative"
            :class="featuresVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
            :style="`transition: all 0.8s cubic-bezier(0.16, 1, 0.3, 1) ${index * 0.1}s`"
          >
            <!-- Card -->
            <div class="relative h-full bg-background-light dark:bg-background-dark rounded-2xl p-8 border border-surface-light dark:border-surface-dark overflow-hidden transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-primary/5">
              <!-- Background Gradient -->
              <div 
                class="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500"
                :class="{
                  'bg-gradient-to-br from-primary/5 to-transparent': feature.color === 'primary',
                  'bg-gradient-to-br from-secondary/5 to-transparent': feature.color === 'secondary',
                  'bg-gradient-to-br from-accent/5 to-transparent': feature.color === 'accent'
                }"
              />
              
              <!-- Content -->
              <div class="relative">
                <!-- Icon -->
                <div 
                  class="w-14 h-14 rounded-xl flex items-center justify-center mb-6 transition-transform duration-300 group-hover:scale-110 group-hover:rotate-3"
                  :class="{
                    'bg-primary/10 text-primary': feature.color === 'primary',
                    'bg-secondary/10 text-secondary': feature.color === 'secondary',
                    'bg-accent/10 text-accent': feature.color === 'accent'
                  }"
                >
                  <svg v-if="feature.color === 'primary'" class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
                  </svg>
                  <svg v-else-if="feature.color === 'secondary'" class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <svg v-else class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                  </svg>
                </div>

                <!-- Title -->
                <h3 class="text-xl font-bold text-text-light dark:text-text-dark mb-3">
                  {{ feature.title }}
                </h3>

                <!-- Description -->
                <p class="text-text-light/70 dark:text-text-dark/70 leading-relaxed">
                  {{ feature.description }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-24 relative overflow-hidden">
      <div class="absolute inset-0 bg-gradient-to-r from-primary/5 via-secondary/5 to-accent/5" />
      <div class="relative mx-auto max-w-4xl px-6 lg:px-8 text-center">
        <h2 class="text-3xl sm:text-4xl font-bold text-text-light dark:text-text-dark mb-6">
          Ready to Get Started?
        </h2>
        <p class="text-lg text-text-light/70 dark:text-text-dark/70 mb-8 max-w-xl mx-auto">
          Join thousands of individuals who have transformed their productivity with Project Life Ledger.
        </p>
        <BaseButton 
          @click="navigateToApp" 
          text="Start Your Journey" 
          variant="primary" 
          size="lg"
          class="transform hover:scale-105 transition-all duration-200 shadow-lg shadow-primary/25"
        />
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-surface-light dark:bg-surface-dark border-t border-surface-light dark:border-surface-dark">
      <div class="mx-auto max-w-7xl px-6 py-12 lg:px-8">
        <div class="flex flex-col md:flex-row items-center justify-between gap-6">
          <!-- Logo -->
          <div class="flex items-center space-x-3">
            <img src="/logo.svg" alt="Project Life Ledger" class="h-8 w-8" />
            <span class="text-lg font-bold text-text-light dark:text-text-dark">
              Project Life Ledger
            </span>
          </div>

          <!-- Links -->
          <div class="flex items-center gap-6">
            <NuxtLink 
              to="/privacy-policy"
              class="text-sm text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark transition-colors"
            >
              Privacy Policy
            </NuxtLink>
            <NuxtLink 
              to="/terms-of-service"
              class="text-sm text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark transition-colors"
            >
              Terms of Service
            </NuxtLink>
          </div>

          <!-- Copyright -->
          <p class="text-sm text-text-light/50 dark:text-text-dark/50">
            © 2024 Project Life Ledger
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Isometric Styles */
.isometric-scene {
  transform-style: preserve-3d;
  transform: rotateX(60deg) rotateZ(-45deg);
  width: 300px;
  height: 300px;
  position: relative;
}

.iso-element {
  position: absolute;
  transform-style: preserve-3d;
}

.iso-task {
  top: 20%;
  left: 10%;
}

.iso-finance {
  top: 40%;
  left: 40%;
}

.iso-learning {
  top: 10%;
  left: 50%;
}

.iso-card {
  width: 140px;
  height: 100px;
  position: relative;
  transform-style: preserve-3d;
}

.iso-card-face {
  position: absolute;
  border-radius: 8px;
}

.iso-card-front {
  width: 140px;
  height: 100px;
  background: rgba(255, 255, 255, 0.9);
  dark: bg-surface-dark;
  border: 1px solid rgba(0, 0, 0, 0.1);
  transform: translateZ(10px);
  padding: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.iso-card-side {
  width: 20px;
  height: 100px;
  background: rgba(0, 0, 0, 0.05);
  transform: rotateY(-90deg) translateZ(10px);
  left: -10px;
}

.iso-card-top {
  width: 140px;
  height: 20px;
  background: rgba(255, 255, 255, 0.7);
  transform: rotateX(90deg) translateZ(10px);
  top: -10px;
}

.iso-stack {
  transform-style: preserve-3d;
}

.iso-flashcard {
  width: 120px;
  height: 80px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  padding: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transform-style: preserve-3d;
}

/* Floating Animation */
@keyframes floating {
  0%, 100% {
    transform: translateY(0) translateZ(0);
  }
  50% {
    transform: translateY(-15px) translateZ(10px);
  }
}

.floating {
  animation: floating 4s ease-in-out infinite;
}

/* Smooth Scroll Behavior */
html {
  scroll-behavior: smooth;
}

/* Dark Mode Support */
.dark .iso-card-front,
.dark .iso-flashcard {
  background: rgba(30, 30, 30, 0.9);
  border-color: rgba(255, 255, 255, 0.1);
}

.dark .iso-card-top {
  background: rgba(50, 50, 50, 0.7);
}

.dark .iso-card-side {
  background: rgba(255, 255, 255, 0.05);
}
</style>
