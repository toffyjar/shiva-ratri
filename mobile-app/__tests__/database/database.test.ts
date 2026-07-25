const mockDb = {
  execAsync: jest.fn().mockResolvedValue(undefined),
  runAsync: jest.fn().mockResolvedValue(undefined),
  getAllAsync: jest.fn().mockResolvedValue([]),
};

jest.mock('expo-sqlite', () => ({
  openDatabaseAsync: jest.fn().mockResolvedValue(mockDb),
}));

import { initDatabase, saveKundli, getAllKundli } from '../../database/database';

describe('database', () => {
  beforeEach(() => {
    mockDb.execAsync.mockClear();
    mockDb.runAsync.mockClear();
    mockDb.getAllAsync.mockClear();
  });

  describe('initDatabase', () => {
    it('creates the kundli table if it does not exist', async () => {
      await initDatabase();

      expect(mockDb.execAsync).toHaveBeenCalledTimes(1);
      const sql = mockDb.execAsync.mock.calls[0][0];
      expect(sql).toContain('CREATE TABLE IF NOT EXISTS kundli');
    });
  });

  describe('saveKundli', () => {
    it('inserts a kundli row with fields in the expected order', async () => {
      await saveKundli({
        rashi: 3,
        planets: 'Su,Mo',
        name: 'Arjun',
        birth_time: '10:30',
        birth_place: 'Delhi',
        gender: 'Male',
      });

      expect(mockDb.runAsync).toHaveBeenCalledTimes(1);
      const [sql, params] = mockDb.runAsync.mock.calls[0];
      expect(sql).toContain('INSERT INTO kundli');
      expect(params).toEqual([3, 'Su,Mo', 'Arjun', '10:30', 'Delhi', 'Male']);
    });
  });

  describe('getAllKundli', () => {
    it('returns rows from getAllAsync', async () => {
      const rows = [
        { id: 1, rashi: 3, planets: 'Su', name: 'Arjun', birth_time: '10:30', birth_place: 'Delhi', gender: 'Male' },
      ];
      mockDb.getAllAsync.mockResolvedValueOnce(rows);

      const result = await getAllKundli();

      expect(mockDb.getAllAsync).toHaveBeenCalledWith('SELECT * FROM kundli');
      expect(result).toEqual(rows);
    });

    it('returns an empty array when there are no rows', async () => {
      mockDb.getAllAsync.mockResolvedValueOnce([]);

      const result = await getAllKundli();

      expect(result).toEqual([]);
    });
  });
});
