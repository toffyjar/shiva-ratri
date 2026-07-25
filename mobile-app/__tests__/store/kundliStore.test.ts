import { useKundliStore } from '../../store/kundliStore';

const initialState = {
  name: '',
  day: '',
  month: '',
  year: '',
  hour: '',
  minute: '',
  birthPlace: '',
  latitude: '',
  longitude: '',
  timeZone: '+00.00',
  isFemale: 'Male',
};

describe('useKundliStore', () => {
  beforeEach(() => {
    useKundliStore.setState(initialState);
  });

  it('has the expected initial state', () => {
    const state = useKundliStore.getState();
    expect(state.name).toBe('');
    expect(state.timeZone).toBe('+00.00');
    expect(state.isFemale).toBe('Male');
  });

  it.each([
    ['setName', 'name', 'Arjun'],
    ['setDay', 'day', '15'],
    ['setMonth', 'month', '08'],
    ['setYear', 'year', '1995'],
    ['setHour', 'hour', '10'],
    ['setMinute', 'minute', '30'],
    ['setBirthPlace', 'birthPlace', 'Delhi, India'],
    ['setLatitude', 'latitude', '28.6139'],
    ['setLongitude', 'longitude', '77.2090'],
    ['setTimeZone', 'timeZone', '+05:30'],
    ['setIsFemale', 'isFemale', 'Female'],
  ] as const)('%s updates only the %s field', (action, field, value) => {
    const before = useKundliStore.getState();
    (before[action] as (v: string) => void)(value);

    const after = useKundliStore.getState();
    expect(after[field]).toBe(value);

    // every other field should be unchanged
    for (const key of Object.keys(initialState) as (keyof typeof initialState)[]) {
      if (key !== field) {
        expect(after[key]).toBe(initialState[key]);
      }
    }
  });

  it('reset restores the initial state after multiple mutations', () => {
    const { setName, setDay, setBirthPlace, setIsFemale, reset } = useKundliStore.getState();
    setName('Arjun');
    setDay('15');
    setBirthPlace('Delhi');
    setIsFemale('Female');

    reset();

    const state = useKundliStore.getState();
    expect(state.name).toBe('');
    expect(state.day).toBe('');
    expect(state.birthPlace).toBe('');
    expect(state.isFemale).toBe('Male');
    expect(state.timeZone).toBe('+00.00');
  });
});
