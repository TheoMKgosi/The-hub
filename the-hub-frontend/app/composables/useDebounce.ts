/**
 * Creates a debounced function that delays invoking func until after wait milliseconds
 * have elapsed since the last time the debounced function was invoked.
 * 
 * @param func - The function to debounce
 * @param wait - The number of milliseconds to delay
 * @param immediate - If true, trigger the function on the leading edge instead of trailing
 * @returns The debounced function with a cancel method
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number = 500,
  immediate: boolean = false
): T & { cancel: () => void } {
  let timeoutId: NodeJS.Timeout | null = null;
  let lastArgs: Parameters<T> | null = null;

  const debounced = function (this: any, ...args: Parameters<T>) {
    lastArgs = args;
    const callNow = immediate && !timeoutId;

    const later = () => {
      timeoutId = null;
      if (!immediate && lastArgs) {
        func.apply(this, lastArgs);
      }
    };

    if (timeoutId) {
      clearTimeout(timeoutId);
    }
    
    timeoutId = setTimeout(later, wait);

    if (callNow) {
      func.apply(this, args);
    }
  } as T & { cancel: () => void };

  // Add cancel method to clear pending execution
  debounced.cancel = () => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      timeoutId = null;
    }
  };

  return debounced;
}
