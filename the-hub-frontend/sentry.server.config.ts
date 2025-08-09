import * as Sentry from "@sentry/nuxt";

if (process.env.NODE_ENV === 'production') {
  Sentry.init({
    dsn: "https://0f93f48bf9daff70a1730cd729955dc0@o4509804910936064.ingest.de.sentry.io/4509804913557584",

    // We recommend adjusting this value in production, or using tracesSampler
    // for finer control
    tracesSampleRate: 1.0,

    // Enable logs to be sent to Sentry
    enableLogs: true,

    // Setting this option to true will print useful information to the console while you're setting up Sentry.
    debug: false,
  });
}
