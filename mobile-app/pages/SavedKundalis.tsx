import * as React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { List, MD3Colors } from 'react-native-paper';

import { savedKundalisPageStyles as styles } from '../custom-styles/savedKundalisPageStyles';
import useKundliStore from '../store/kundliStore';
import { useAuthStore } from '../store/authStore';
import { useSavedKundalisStore } from '../store/savedKundalisStore';

type SavedKundalisProps = {
  onSelectKundali?: () => void;
};

export default function SavedKundalis({ onSelectKundali }: SavedKundalisProps) {
  const loadKundli = useKundliStore((state) => state.loadKundli);
  const token = useAuthStore((state) => state.token);
  const kundalis = useSavedKundalisStore((state) => state.kundalis);
  const fetchKundalis = useSavedKundalisStore((state) => state.fetchKundalis);

  React.useEffect(() => {
    fetchKundalis(token);
  }, [token]);

  return (
    <View style={styles.container}>
      <List.Section>
        {kundalis.map((kundali) => (
          <List.Item
            key={kundali.id}
            title={kundali.name}
            left={() => <List.Icon icon="file" />}
            onPress={() => {
              loadKundli(kundali);
              onSelectKundali?.();
            }}
          />
        ))}
      </List.Section>
    </View>
  );
}

