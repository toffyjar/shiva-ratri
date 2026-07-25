import { formatTimeZoneForApi, isChartData, formatChartData } from '../../utils/chartPage';

describe('formatTimeZoneForApi', () => {
  it('formats a positive fractional offset', () => {
    expect(formatTimeZoneForApi('5.5')).toBe('+05:30');
  });

  it('formats a negative fractional offset', () => {
    expect(formatTimeZoneForApi('-5.5')).toBe('-05:30');
  });

  it('formats a whole-hour offset', () => {
    expect(formatTimeZoneForApi('9')).toBe('+09:00');
  });

  it('formats zero as positive', () => {
    expect(formatTimeZoneForApi('0')).toBe('+00:00');
  });

  it('falls back to zero for a non-numeric string', () => {
    expect(formatTimeZoneForApi('abc')).toBe('+00:00');
  });

  it('falls back to zero for an empty string', () => {
    expect(formatTimeZoneForApi('')).toBe('+00:00');
  });

  it('pads single-digit hours and minutes', () => {
    expect(formatTimeZoneForApi('1.25')).toBe('+01:15');
  });

  it('handles a negative whole-hour offset', () => {
    expect(formatTimeZoneForApi('-12')).toBe('-12:00');
  });
});

describe('isChartData', () => {
  it('returns true for a well-formed chart data object', () => {
    const value = {
      1: { rashi: '1', planets: ['Su', 'Mo'] },
      2: { rashi: '2', planets: [] },
    };
    expect(isChartData(value)).toBe(true);
  });

  it('returns true for an empty object (vacuously)', () => {
    expect(isChartData({})).toBe(true);
  });

  it('returns false for null', () => {
    expect(isChartData(null)).toBe(false);
  });

  it('returns false for undefined', () => {
    expect(isChartData(undefined)).toBe(false);
  });

  it('returns false for a non-object primitive', () => {
    expect(isChartData('not an object')).toBe(false);
    expect(isChartData(42)).toBe(false);
  });

  it('returns false when an entry is missing "planets"', () => {
    expect(isChartData({ 1: { rashi: '1' } })).toBe(false);
  });

  it('returns false when an entry is missing "rashi"', () => {
    expect(isChartData({ 1: { planets: [] } })).toBe(false);
  });

  it('returns false when "planets" is not an array', () => {
    expect(isChartData({ 1: { rashi: '1', planets: 'Su' } })).toBe(false);
  });

  it('returns false when an entry value is not an object', () => {
    expect(isChartData({ 1: 'not-an-object' })).toBe(false);
  });
});

describe('formatChartData', () => {
  it('builds house entries with rashi and planet lists', () => {
    const chart = {
      lagna: { Lg: { rashi: '3' } },
      houses: {
        '1': { graha: { Su: {}, Mo: {} } },
        '2': { graha: {} },
      },
    };

    const result = formatChartData(chart);

    expect(result[1]).toEqual({ rashi: '3', planets: ['Su', 'Mo'] });
    expect(result[2]).toEqual({ rashi: '4', planets: [] });
  });

  it('wraps rashi numbers greater than 12 back around to 1', () => {
    const chart = {
      lagna: { Lg: { rashi: '11' } },
      houses: {
        '3': { graha: {} },
      },
    };

    const result = formatChartData(chart);

    // startAsc(11) + (house(3) - 1) = 13 -> wraps to 1
    expect(result[3]).toEqual({ rashi: '1', planets: [] });
  });

  it('returns an empty object when there are no houses', () => {
    const chart = {
      lagna: { Lg: { rashi: '1' } },
      houses: {},
    };

    expect(formatChartData(chart)).toEqual({});
  });
});
