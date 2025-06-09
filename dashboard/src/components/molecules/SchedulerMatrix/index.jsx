import Button from '@/components/atom/Button';
import React, { useState, useMemo, useCallback } from 'react';

// Immutable configuration
const CONFIG = Object.freeze({
  DAYS_IN_WEEK: 7,
  HOURS_IN_DAY: 24,
  MINUTES_IN_HOUR: 60,
  SLOT_DURATIONS: Object.freeze({
    '15min': 15,
    '30min': 30,
    '60min': 60,
  }),
  SLOT_STATES: Object.freeze({
    INACTIVE: 0,
    ACTIVE: 1,
  }),
});

// Immutable day definitions
const DAYS = Object.freeze([
  { id: 0, label: 'MON', fullName: 'Monday' },
  { id: 1, label: 'TUE', fullName: 'Tuesday' },
  { id: 2, label: 'WED', fullName: 'Wednesday' },
  { id: 3, label: 'THU', fullName: 'Thursday' },
  { id: 4, label: 'FRI', fullName: 'Friday' },
  { id: 5, label: 'SAT', fullName: 'Saturday' },
  { id: 6, label: 'SUN', fullName: 'Sunday' },
]);

// Memoized time slots
const TIME_SLOTS = Object.freeze(
  Array.from({ length: CONFIG.HOURS_IN_DAY }, (_, i) => ({
    id: i,
    label: i.toString().padStart(2, '0'),
    hour: i,
  })),
);

// Pure utility functions
const getSlotsPerHour = (resolution) => CONFIG.MINUTES_IN_HOUR / CONFIG.SLOT_DURATIONS[resolution];

const getSlotIndex = (hourIndex, subSlotIndex, resolution) =>
  hourIndex * getSlotsPerHour(resolution) + subSlotIndex;

const createEmptySchedule = (resolution = '30min') => {
  const slotDuration = CONFIG.SLOT_DURATIONS[resolution];
  const slotsPerDay = (CONFIG.HOURS_IN_DAY * CONFIG.MINUTES_IN_HOUR) / slotDuration;

  return Array.from({ length: CONFIG.DAYS_IN_WEEK }, () =>
    Array(slotsPerDay).fill(CONFIG.SLOT_STATES.INACTIVE),
  );
};

const getSlotInfo = (dayIndex, slotIndex, resolution) => {
  const slotsPerHour = getSlotsPerHour(resolution);
  const hour = Math.floor(slotIndex / slotsPerHour);
  const subSlot = slotIndex % slotsPerHour;
  const minutes = subSlot * CONFIG.SLOT_DURATIONS[resolution];

  return {
    day: DAYS[dayIndex],
    hour,
    subSlot,
    minutes,
    timeString: `${hour.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`,
  };
};

// Sample data factory with functional approach
const createSampleSchedule = () => {
  const schedule = createEmptySchedule('15min');
  const { ACTIVE, INACTIVE } = CONFIG.SLOT_STATES;

  const createPattern = (patterns) => {
    const result = Array(96).fill(INACTIVE);
    patterns.forEach(({ start, end, state }) => {
      for (let i = start; i < end; i++) {
        result[i] = state;
      }
    });
    return result;
  };

  const workdayPattern = createPattern([
    { start: 0, end: 8, state: ACTIVE }, // Early morning
    { start: 32, end: 72, state: ACTIVE }, // Working hours
    { start: 48, end: 52, state: INACTIVE }, // Lunch break
    { start: 72, end: 80, state: ACTIVE }, // Evening
    { start: 88, end: 96, state: ACTIVE }, // Night
  ]);

  const weekendPattern = createPattern([
    { start: 0, end: 36, state: ACTIVE }, // Morning
    { start: 36, end: 44, state: INACTIVE }, // Break
    { start: 44, end: 80, state: ACTIVE }, // Afternoon
    { start: 80, end: 88, state: INACTIVE }, // Break
    { start: 88, end: 96, state: ACTIVE }, // Night
  ]);

  // Apply patterns efficiently
  [0, 1, 2, 3].forEach((day) => {
    schedule[day] = [...workdayPattern];
  });

  // Friday with extended evening activity
  schedule[4] = workdayPattern.map((val, i) => (i >= 72 ? ACTIVE : val));

  [5, 6].forEach((day) => {
    schedule[day] = [...weekendPattern];
  });

  return schedule;
};

// Immutable schedule manager using functional approach
class ImmutableScheduleManager {
  constructor(schedule, resolution = '15min') {
    this.schedule = schedule;
    this.resolution = resolution;
  }

  getSlotState(dayIndex, hourIndex, subSlotIndex) {
    const slotIndex = getSlotIndex(hourIndex, subSlotIndex, this.resolution);
    return this.schedule[dayIndex]?.[slotIndex] ?? CONFIG.SLOT_STATES.INACTIVE;
  }

  _createNewSchedule(updater) {
    const newSchedule = updater(this.schedule.map((day) => [...day]));
    return new ImmutableScheduleManager(newSchedule, this.resolution);
  }

  setSlotState(dayIndex, hourIndex, subSlotIndex, state) {
    const slotIndex = getSlotIndex(hourIndex, subSlotIndex, this.resolution);

    return this._createNewSchedule((schedule) => {
      schedule[dayIndex][slotIndex] = state;
      return schedule;
    });
  }

  toggleSlot(dayIndex, hourIndex, subSlotIndex) {
    const currentState = this.getSlotState(dayIndex, hourIndex, subSlotIndex);
    const newState =
      currentState === CONFIG.SLOT_STATES.ACTIVE
        ? CONFIG.SLOT_STATES.INACTIVE
        : CONFIG.SLOT_STATES.ACTIVE;

    return this.setSlotState(dayIndex, hourIndex, subSlotIndex, newState);
  }

  toggleDay(dayIndex) {
    const daySchedule = this.schedule[dayIndex];
    const allInactive = daySchedule.every((slot) => slot === CONFIG.SLOT_STATES.INACTIVE);
    const newState = allInactive ? CONFIG.SLOT_STATES.ACTIVE : CONFIG.SLOT_STATES.INACTIVE;

    return this._createNewSchedule((schedule) => {
      schedule[dayIndex] = Array(daySchedule.length).fill(newState);
      return schedule;
    });
  }

  toggleHour(hourIndex) {
    const slotsPerHour = getSlotsPerHour(this.resolution);

    // Check if all slots in this hour are inactive across all days
    const allHourSlots = this.schedule.flatMap((day) =>
      Array.from({ length: slotsPerHour }, (_, subSlot) => {
        const slotIndex = getSlotIndex(hourIndex, subSlot, this.resolution);
        return day[slotIndex];
      }),
    );

    const allInactive = allHourSlots.every((slot) => slot === CONFIG.SLOT_STATES.INACTIVE);
    const newState = allInactive ? CONFIG.SLOT_STATES.ACTIVE : CONFIG.SLOT_STATES.INACTIVE;

    return this._createNewSchedule((schedule) => {
      schedule.forEach((day) => {
        for (let subSlot = 0; subSlot < slotsPerHour; subSlot++) {
          const slotIndex = getSlotIndex(hourIndex, subSlot, this.resolution);
          day[slotIndex] = newState;
        }
      });
      return schedule;
    });
  }

  getActiveSlots() {
    return this.schedule.flatMap((daySchedule, dayIndex) =>
      daySchedule
        .map((slotState, slotIndex) => ({ slotState, slotIndex, dayIndex }))
        .filter(({ slotState }) => slotState === CONFIG.SLOT_STATES.ACTIVE)
        .map(({ slotIndex, dayIndex }) => {
          const slotInfo = getSlotInfo(dayIndex, slotIndex, this.resolution);
          return {
            ...slotInfo,
            day: slotInfo.day.label,
            time: slotInfo.timeString,
            dayIndex,
            slotIndex,
          };
        }),
    );
  }

  getScheduleStats() {
    const totalSlots = this.schedule.flat().length;
    const activeSlots = this.schedule
      .flat()
      .filter((slot) => slot === CONFIG.SLOT_STATES.ACTIVE).length;
    const inactiveSlots = totalSlots - activeSlots;

    const slotDurationMinutes = CONFIG.SLOT_DURATIONS[this.resolution];
    const activeHours = (activeSlots * slotDurationMinutes) / 60;
    const inactiveHours = (inactiveSlots * slotDurationMinutes) / 60;
    const inactivePercentage = ((inactiveSlots / totalSlots) * 100).toFixed(1);

    return {
      activeHours: Math.round(activeHours * 10) / 10,
      inactiveHours: Math.round(inactiveHours * 10) / 10,
      inactivePercentage: parseFloat(inactivePercentage),
    };
  }

  export() {
    return {
      schedule: this.schedule.map((day) => [...day]),
      resolution: this.resolution,
      metadata: {
        totalDays: CONFIG.DAYS_IN_WEEK,
        slotsPerDay: this.schedule[0]?.length ?? 0,
        slotsPerHour: getSlotsPerHour(this.resolution),
        slotDuration: CONFIG.SLOT_DURATIONS[this.resolution],
      },
    };
  }
}

const SlotCell = React.memo(
  ({ dayIndex, hourIndex, scheduleManager, timeResolution, onSlotToggle }) => {
    const slotsPerHour = getSlotsPerHour(timeResolution);
    const isMultiSlot = slotsPerHour > 1;

    const subSlots = useMemo(
      () => Array.from({ length: slotsPerHour }, (_, i) => i),
      [slotsPerHour],
    );

    if (isMultiSlot) {
      return (
        <div className="w-[30px] h-8 m-1 flex mx-1">
          {subSlots.map((subSlotIndex) => {
            const slotState = scheduleManager.getSlotState(dayIndex, hourIndex, subSlotIndex);
            const bgColor = slotState === CONFIG.SLOT_STATES.ACTIVE ? 'bg-green-500' : 'bg-red-500';
            const minutes = subSlotIndex * CONFIG.SLOT_DURATIONS[timeResolution];

            return (
              <div
                key={subSlotIndex}
                className={`w-full h-full flex-1 mx-[1px] rounded-lg ${bgColor} cursor-pointer hover:opacity-80 transition-opacity`}
                onClick={() => onSlotToggle(dayIndex, hourIndex, subSlotIndex)}
                title={`${DAYS[dayIndex].fullName} ${TIME_SLOTS[hourIndex].label}:${minutes.toString().padStart(2, '0')}`}
              />
            );
          })}
        </div>
      );
    }

    const slotState = scheduleManager.getSlotState(dayIndex, hourIndex, 0);
    const bgColor = slotState === CONFIG.SLOT_STATES.ACTIVE ? 'bg-green-500' : 'bg-red-500';

    return (
      <div className="w-[30px] h-8 m-1 flex mx-1">
        <div
          className={`w-full h-full cursor-pointer hover:opacity-80 rounded-lg transition-opacity ${bgColor}`}
          onClick={() => onSlotToggle(dayIndex, hourIndex, 0)}
          title={`${DAYS[dayIndex].fullName} ${TIME_SLOTS[hourIndex].label}:00`}
        />
      </div>
    );
  },
);

SlotCell.displayName = 'SlotCell';

export default function SchedulerMatrix() {
  const [timeResolution, setTimeResolution] = useState('15min');

  // Memoize initial schedule manager
  const initialScheduleManager = useMemo(
    () => new ImmutableScheduleManager(createSampleSchedule(), '15min'),
    [],
  );

  const [scheduleManager, setScheduleManager] = useState(initialScheduleManager);

  // Memoized handlers
  const handleSlotToggle = useCallback((dayIndex, hourIndex, subSlotIndex) => {
    setScheduleManager((prev) => prev.toggleSlot(dayIndex, hourIndex, subSlotIndex));
  }, []);

  const handleDayToggle = useCallback((dayIndex) => {
    setScheduleManager((prev) => prev.toggleDay(dayIndex));
  }, []);

  const handleHourToggle = useCallback((hourIndex) => {
    setScheduleManager((prev) => prev.toggleHour(hourIndex));
  }, []);

  const handleApply = useCallback(() => {
    const exportedData = scheduleManager.export();
    const activeSlots = scheduleManager.getActiveSlots();
    const stats = scheduleManager.getScheduleStats();
    console.log(stats);

    console.log('Schedule Export:', exportedData);
    console.log('Active Slots:', activeSlots);
    console.log('Raw Schedule Data:', exportedData.schedule);
  }, [scheduleManager]);

  const handleResolutionChange = useCallback((newResolution) => {
    const newSchedule = createEmptySchedule(newResolution);
    const newManager = new ImmutableScheduleManager(newSchedule, newResolution);
    setTimeResolution(newResolution);
    setScheduleManager(newManager);
  }, []);

  // Memoized active slots count
  const activeSlotCount = useMemo(() => scheduleManager.getActiveSlots().length, [scheduleManager]);
  const stats = useMemo(() => scheduleManager.getScheduleStats().length, [scheduleManager]);

  // Memoized resolution options
  const resolutionOptions = useMemo(() => Object.keys(CONFIG.SLOT_DURATIONS), []);

  return (
    <div className="mx-auto p-4 bg-white">
      {/* Controls */}
      <div className="mb-4 flex gap-4 items-center">
        <div className="flex gap-2">
          <label className="text-sm font-medium">Resolution:</label>
          {resolutionOptions.map((resolution) => (
            <button
              key={resolution}
              onClick={() => handleResolutionChange(resolution)}
              className={`px-3 py-1 text-xs rounded transition-colors ${
                timeResolution === resolution
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              {resolution}
            </button>
          ))}
        </div>
        <div className="text-sm text-gray-600">{activeSlotCount} active slots</div>
      </div>

      <div className="flex flex-wrap items-center justify-between mb-4">
        <div> Estimated Savings </div>
        <div>Active Hours: Inactive Hours: {}</div>
        <div className="flex items-center space-x-6 mb-2">
          <div className="flex items-center text-gray-700">
            <div className="w-4 h-4 bg-green-500 mr-2"></div>
            <span>RUNNING</span>
          </div>
          <div className="flex items-center text-gray-700">
            <div className="w-4 h-4 bg-red-500 mr-2"></div>
            <span>STOPPED</span>
          </div>
        </div>

        {/* <Button onClick={handleReset}>Reset to Initial Data</Button> */}
      </div>

      {/* Matrix Grid */}
      <div className="overflow-auto">
        <div className="xs:min-w-max lg:w-full">
          {/* Time Slot Headers */}
          <div className="flex">
            <div className="w-10 flex-shrink-0" />
            {TIME_SLOTS.map((slot) => (
              <div
                key={slot.id}
                onClick={() => handleHourToggle(slot.id)}
                className="w-[30px] text-center hover:bg-primary-200 cursor-pointer border text-xs font-semibold border-primary-200 bg-primary-100 py-1 mx-1 my-2 rounded-lg text-gray-700 transition-colors"
                title={`Toggle hour ${slot.label}:00 for all days`}
              >
                {slot.label}
              </div>
            ))}
          </div>

          {/* Days and Time Slots */}
          {DAYS.map((day, dayIndex) => (
            <div key={day.id} className="flex justify-start items-center text-gray-700">
              {/* Day Label */}
              <div
                onClick={() => handleDayToggle(dayIndex)}
                className="w-10 flex-shrink-0 cursor-pointer hover:bg-primary-200 font-semibold flex items-center text-xs rounded-lg justify-center px-2 py-2 bg-primary-100 border border-primary-200 text-gray-700 transition-colors"
                title={`Toggle all slots for ${day.fullName}`}
              >
                {day.label}
              </div>

              {/* Time Slots for this day */}
              {TIME_SLOTS.map((timeSlot) => (
                <SlotCell
                  key={`${dayIndex}-${timeSlot.id}`}
                  dayIndex={dayIndex}
                  hourIndex={timeSlot.id}
                  scheduleManager={scheduleManager}
                  timeResolution={timeResolution}
                  onSlotToggle={handleSlotToggle}
                />
              ))}
            </div>
          ))}
        </div>
      </div>

      <div className="flex items-center justify-end my-4">
        <Button onClick={handleApply}>Apply</Button>
      </div>
    </div>
  );
}
