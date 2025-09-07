import * as Sentry from "@sentry/nuxt";

if (process.env.NODE_ENV === 'production') {
  Sentry.init({
    dsn: useRuntimeConfig().public.sentry.dsn,

    // Reduced sample rate for better performance
    tracesSampleRate: 0.1,

    // This sets the sample rate to be 10%. You may want this to be 100% while
    // in development and sample at a lower rate in production
    replaysSessionSampleRate: 0.1,

    // If the entire session is not sampled, use the below sample rate to sample
    // sessions when an error occurs.
    replaysOnErrorSampleRate: 1.0,

    // If you don't want to use Session Replay, just remove the line below:
    integrations: [Sentry.replayIntegration()],

    // Disable excessive logging to reduce noise
    enableLogs: false,

    // Setting this option to true will print useful information to the console while you're setting up Sentry.
    debug: false,
  });

}
