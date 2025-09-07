import * as Sentry from "@sentry/nuxt";

if (process.env.NODE_ENV === 'production') {
  Sentry.init({
    dsn: useRuntimeConfig().public.sentry.dsn,

    // Reduced sample rate for better performance
    tracesSampleRate: 0.1,

    // Disable excessive logging to reduce noise
    enableLogs: false,

    // Setting this option to true will print useful information to the console while you're setting up Sentry.
    debug: false,
  });
}
