import { create } from 'zustand';
import { authHttpService } from '../services/httpService';
import { SavedKundaliRecord } from '../utils/chartPage';

export type SavedKundali = SavedKundaliRecord & { id: number };

type SavedKundalisState = {
  kundalis: SavedKundali[];
  isLoading: boolean;
  fetchKundalis: (token?: string | null) => Promise<void>;
  addKundali: (kundali: SavedKundali) => void;
};

export const useSavedKundalisStore = create<SavedKundalisState>((set) => ({
  kundalis: [],
  isLoading: false,

  fetchKundalis: async (token) => {
    set({ isLoading: true });

    try {
      const response = await authHttpService.get<{ data: SavedKundali[] }>(
        '/list-kundali',
        undefined,
        token ? { Authorization: `Bearer ${token}` } : undefined
      );

      if (response.ok && response.data?.data) {
        set({ kundalis: response.data.data });
      }
    } catch (error) {
      console.log('Error fetching kundalis:', error);
    } finally {
      set({ isLoading: false });
    }
  },

  addKundali: (kundali) => set((state) => ({ kundalis: [...state.kundalis, kundali] })),
}));

export default useSavedKundalisStore;
