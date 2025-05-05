import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import utc from 'dayjs/plugin/utc';
import isBetween from 'dayjs/plugin/isBetween';
import duration from 'dayjs/plugin/duration';

dayjs.extend(relativeTime);
dayjs.extend(utc);
dayjs.extend(isBetween); // Extend dayjs with the isBetween plugin
dayjs.extend(duration);

const getOrdinalSuffix = (day) => {
  if (day > 3 && day < 21) return 'th'; // Covers 11th, 12th, 13th, etc.
  const suffixes = {
    1: 'st',
    2: 'nd',
    3: 'rd',
  };
  return suffixes[day % 10] || 'th';
};

export const formatTime = (time, type) => {
  const utcTime = dayjs.utc(time);
  const localTime = utcTime.local();
  const timeDifference = dayjs().diff(utcTime, 'second');

  if (timeDifference < 5) {
    return 'Just now';
  } else if (timeDifference < 60) {
    return `${timeDifference} seconds ago`;
  } else if (timeDifference < 60 * 60) {
    return `${Math.floor(timeDifference / 60)} minutes ago`;
  } else if (timeDifference < 60 * 60 * 24) {
    return `${Math.floor(timeDifference / (60 * 60))} hours ago`;
  } else {
    const startOfToday = dayjs().startOf('day');
    const startOfYesterday = startOfToday.subtract(1, 'day');

    if (localTime.isBetween(startOfYesterday, startOfToday)) {
      return `Yesterday, ${localTime.format('HH:mm')}`;
    } else {
      const day = localTime.date();
      const month = localTime.format('MMMM');
      const year = localTime.year();
      const formattedTime = localTime.format('HH:mm');
      const ordinalSuffix = getOrdinalSuffix(day);
      const formattedDate = `${day}${ordinalSuffix} ${month} ${year}`;
      if (type === 'table') {
        return `on ${formattedDate}, ${formattedTime}`;
      } else {
        return `${formattedDate}, ${formattedTime}`;
      }
    }
  }
};

export const calculateAge = (creationTimestamp) => {
  const now = dayjs();
  const creationTime = dayjs(creationTimestamp);
  const diffInMs = now.diff(creationTime);

  // Calculate days, hours, minutes
  const days = Math.floor(diffInMs / (1000 * 60 * 60 * 24));
  const hours = Math.floor((diffInMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  const minutes = Math.floor((diffInMs % (1000 * 60 * 60)) / (1000 * 60));

  let result = '';
  if (days > 0) result += `${days}d`;
  if (hours > 0) result += `${result ? ' ' : ''}${hours}h`;
  if (minutes > 0) result += `${result ? ' ' : ''}${minutes}m`;

  return result || '0m'; // Default to "0m" if the difference is negligible
};

export function ParseJSON(str) {
  try {
    if (typeof str !== 'string') {
      throw new Error('input is not string');
    }
    if (typeof str === 'string' && !isNaN(str)) {
      throw new Error('Invalid JSON: Numeric string');
    }
    const result = JSON.parse(str);
    return result;
  } catch {
    return {};
  }
}
