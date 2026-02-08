// src/composables/useDate.js
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

// Extend it here once
dayjs.extend(relativeTime);

export function useDate() {
  
  const fromNow = (date) => {
    if (!date) return '';
    return dayjs(date).fromNow();
  };

  const format = (date, template = 'DD MMM YYYY') => {
    if (!date) return '';
    return dayjs(date).format(template);
  };

  return {
    fromNow,
    format,
    dayjs // Expose the instance if you need more complex logic
  };
}
